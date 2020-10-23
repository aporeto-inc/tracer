package monitoring

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"go.aporeto.io/elemental"
	"go.aporeto.io/gaia"
	"go.uber.org/zap"
)

// APIErrors represent a list of API errors
type APIErrors []APIError

// APIError repesent an API error
type APIError struct {
	Code      int
	Service   string
	Identity  string
	Operation string
	Method    string
	URL       string
	Count     int
	Traces    []string
}

// Hash return a unique identifier for an error
// based on url code and operation
func (a APIError) Hash() uint32 {

	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%d", a.Code) + a.URL + a.Method)) // nolint
	return h.Sum32()
}

// ByCount implements sort.Interface based on the count
type ByCount APIErrors

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Less(i, j int) bool { return a[i].Count < a[j].Count }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// GetAPIErrors retrieve the errors metrics from prometheus as APiErrors
func (m Client) GetAPIErrors(since time.Duration, at time.Time) (APIErrors, error) {

	// query the errors
	errRes, err := m.queryPrometheus(fmt.Sprintf("sum(delta(http_requests_total{code!~'0|500'}[%ds])) by (service,code,method,url) >0", int(since.Seconds())), at)
	if err != nil {
		return nil, err
	}

	// query the 500
	panicRes, err := m.queryPrometheus(fmt.Sprintf("count((http_errors_5xx_total{code='500'} > 0 unless http_errors_5xx_total{code='500'} offset %ds) or ((http_errors_5xx_total{code='500'} - http_errors_5xx_total{code='500'} offset %ds) >0)) by (service,code,method,url) >0", int(since.Seconds()), int(since.Seconds())), at)
	if err != nil {
		return nil, err
	}

	return append(errRes, panicRes...), nil

}

func (m Client) queryPrometheus(query string, at time.Time) (APIErrors, error) {

	promProxy, err := url.Parse("api/datasources/proxy/1")
	if err != nil {
		panic(err)
	}

	client, err := api.NewClient(api.Config{Address: m.url.ResolveReference(promProxy).String(), RoundTripper: m.client.Transport})
	if err != nil {
		return nil, fmt.Errorf("unable to create new prometheus client: %w", err)
	}

	v1api := v1.NewAPI(client)
	tctx, tcancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer tcancel()

	result, warnings, err := v1api.Query(tctx, query, at)

	if err != nil {
		return nil, fmt.Errorf("unable to query prometheus: %w", err)
	}

	if len(warnings) > 0 {
		zap.L().Warn("Warning while querying Prometheus", zap.Strings("warnings", warnings))
	}
	results := parseMetrics(result)

	zap.L().Debug("Quering prometheus", zap.String("query", query), zap.Int("results", len(results)))

	return results, nil
}

func parseMetrics(result model.Value) []APIError {

	res := []APIError{}

	keys := []model.LabelName{"code", "method", "url", "service"}

	for _, v := range result.(model.Vector) {

		// Sanitize results
		if !func() bool {
			for _, key := range keys {
				if _, ok := v.Metric[key]; !ok {
					zap.L().Debug("Unable to parse metrics, label not found", zap.String("label", string(key)), zap.Reflect("metric", v.Metric))
					return false
				}
			}
			return true
		}() {
			continue
		}

		code, err := strconv.Atoi(string(v.Metric["code"]))
		if err != nil {
			zap.L().Error("Unable to parse metrics, code is not an integer", zap.String("code", string(v.Metric["code"])))
			continue
		}

		identity, operation, err := extractIdentityFrom(string(v.Metric["url"]), string(v.Metric["method"]))
		if err != nil {
			zap.L().Error("Unable extract identity from url", zap.Error(err))
		}

		res = append(res, APIError{
			Code:      code,
			Identity:  identity.Name,
			Operation: string(operation),
			Service:   string(v.Metric["service"]),
			Method:    string(v.Metric["method"]),
			URL:       string(v.Metric["url"]),
			Count:     int(v.Value),
		})
	}

	return res
}

// extractIdentityFrom extract Identity and Operation from url and method
func extractIdentityFrom(url, method string) (identity elemental.Identity, operation elemental.Operation, err error) {

	manager := gaia.Manager()

	components := strings.Split(url, "/")

	// We remove the first element as it's always empty
	components = append(components[:0], components[1:]...)

	switch len(components) {
	case 1, 2:
		identity = manager.IdentityFromCategory(components[0])
	case 3:
		identity = manager.IdentityFromCategory(components[2])
	default:
		return identity, operation, fmt.Errorf("unable to decode url parts")
	}

	switch method {
	case http.MethodDelete:
		operation = elemental.OperationDelete

	case http.MethodGet:
		if len(components) == 1 || len(components) == 3 {
			operation = elemental.OperationRetrieveMany
		} else {
			operation = elemental.OperationRetrieve
		}

	case http.MethodHead:
		operation = elemental.OperationInfo

	case http.MethodPatch:
		operation = elemental.OperationPatch

	case http.MethodPost:
		operation = elemental.OperationCreate

	case http.MethodPut:
		operation = elemental.OperationUpdate

	}

	return identity, operation, nil

}
