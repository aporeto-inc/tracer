package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aporeto-inc/tracer/internal/monitoring"
)

// Filter is meant to filter APIErrors given filers
func Filter(codes string, services, urls []string, results monitoring.APIErrors) (monitoring.APIErrors, error) {

	serviceFilter := make(map[string]struct{})
	for _, s := range services {
		serviceFilter[s] = struct{}{}
	}

	urlsFilter := make(map[string]struct{})
	for _, s := range urls {
		urlsFilter[s] = struct{}{}
	}

	// Parse code range if any
	codeFilter := make(map[int]struct{})

	// Split parts
	for _, part := range strings.Split(codes, ",") {

		part := strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// is it a range
		bounds := strings.Split(part, "-")

		switch {
		case len(bounds) > 2:
			return nil, fmt.Errorf("--code failed to parse: Bad code range: %s, must be in the form X-Y", part)
		case len(bounds) == 2:
			lb, err := strconv.Atoi(bounds[0])
			if err != nil {
				return nil, fmt.Errorf("--code failed to parse: Unable to convert lower bound to integer: %s", bounds[0])
			}
			ub, err := strconv.Atoi(bounds[1])
			if err != nil {
				return nil, fmt.Errorf("--code failed to parse: Unable to convert upper bound to integer: %s", bounds[1])
			}

			if lb > ub {
				return nil, fmt.Errorf("--code failed to parse: Invalid range, the lower bound: %s is greater than upper bound: %s", bounds[0], bounds[1])
			}

			for i := lb; i <= ub; i++ {
				codeFilter[i] = struct{}{}
			}

		default:
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("--code failed to parse: Unable to convert entry to integer: %s", part)
			}
			codeFilter[p] = struct{}{}
		}
	}

	// Create a hashed list of everything
	toFilter := make(map[uint32]monitoring.APIError)
	for _, result := range results {
		toFilter[result.Hash()] = result
	}

	toRemove := make(map[uint32]struct{})

	// Remove the codes that are not matching
	if len(codeFilter) > 0 {
		for _, result := range results {
			if _, ok := codeFilter[result.Code]; !ok {
				toRemove[result.Hash()] = struct{}{}
			}
		}
	}

	for h := range toRemove {
		delete(toFilter, h)
	}

	// If we have no further filter return the list
	if len(serviceFilter) == 0 && len(urlsFilter) == 0 {
		return func() monitoring.APIErrors {
			out := monitoring.APIErrors{}
			for _, i := range toFilter {
				out = append(out, i)
			}
			return out
		}(), nil
	}

	// Keep only the filters that are mathcing
	filtered := monitoring.APIErrors{}
	for _, result := range toFilter {
		if _, ok := serviceFilter[result.Service]; ok {
			filtered = append(filtered, result)
			continue
		}
		if _, ok := urlsFilter[result.URL]; ok {
			filtered = append(filtered, result)
		}
	}

	return filtered, nil
}
