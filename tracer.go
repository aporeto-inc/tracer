// Package main provides the main definition of tracer, a CLI utility for interacting with our observability stack.
package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aporeto-inc/tracer/internal/configuration"
	"github.com/aporeto-inc/tracer/internal/monitoring"
	"github.com/aporeto-inc/tracer/internal/profiles"
	"github.com/aporeto-inc/tracer/internal/utils"
	"go.uber.org/zap"
)

// Start starts the service
func main() {

	cfg := configuration.NewConfiguration()
	datasource := profiles.NewProfile(cfg)

	if cfg.Open != "" {
		monitoring.OpenTrace(datasource.MonitoringURL, datasource.TracesDataSourceName, cfg.Open)
	}

	var err error
	var from, to time.Time

	from, to, since, err := utils.ParseTime(cfg.From, cfg.To, cfg.Since)
	if err != nil {
		zap.L().Fatal("Unable to parse time", zap.Error(err))
	}

	// Create monitoring client
	c, err := monitoring.NewClient(datasource)
	if err != nil {
		zap.L().Fatal("Unable to create monitoring client", zap.Error(err))
	}

	// Check connection
	if err = c.Ping(); err != nil {
		zap.L().Fatal("Unable to connect to monitoring", zap.Error(err))
	}

	// Show log if asked
	if cfg.Log || cfg.LogFilter != "" {

		quiet := true
		if cfg.LogLevel == "debug" {
			quiet = false
		}

		if err := c.GetLogs(datasource.LogsIndex, from, to, cfg.Services, cfg.LogConf, quiet); err != nil {
			zap.L().Fatal("Unable to get logs", zap.Error(err))
		}

	} else {

		var results monitoring.APIErrors
		// Get the metrics
		results, err = c.GetAPIErrors(datasource.MetricsIndex, since, to)
		if err != nil {
			zap.L().Fatal("Unable to query prometheus", zap.Error(err))
		}

		// Filter
		results, err = utils.Filter(cfg.Codes, cfg.Services, cfg.URLS, results)
		if err != nil {
			zap.L().Fatal("Failed to parse filters", zap.Error(err))
		}

		// Sort by counts
		sort.Sort(monitoring.ByCount(results))

		// Get the traces
		var wg sync.WaitGroup
		wg.Add(len(results))

		for i := range results {
			go func(index int) {
				defer wg.Done()

				params := monitoring.TracingQueryParameters{
					Start:       from.UnixNano() / 1000,
					End:         to.UnixNano() / 1000,
					Limit:       cfg.Limit,
					MinDuration: cfg.MinDuration,
					Service:     strings.Split(results[index].Service, "-")[0],
					Tags: map[string]string{
						"status.code":   fmt.Sprintf("%d", results[index].Code),
						"req.identity":  results[index].Identity,
						"req.operation": results[index].Operation,
					},
				}

				if cfg.OnlyError {
					params.Tags["error"] = "true"
				}

				if cfg.Namespace != "" {
					params.Tags["req.namespace"] = cfg.Namespace
				}

				traceResults, err := c.GetTraceIDs(datasource.TracesIndex, params)
				if err != nil {
					zap.L().Error("Failed to retrieve traces for error", zap.Error(err))
					return
				}

				results[index].Traces = traceResults
			}(i)
		}

		wg.Wait()

		// If we have a trace filter remove the entries without traces
		if cfg.OnlyError || cfg.MinDuration.String() != "0s" || cfg.Namespace != "" {
			results = func() monitoring.APIErrors {
				res := monitoring.APIErrors{}
				for _, item := range results {
					if len(item.Traces) != 0 {
						item.Count = len(item.Traces)
						res = append(res, item)
					}
				}
				return res
			}()
		}

		// Display
		if len(results) > 0 {

			fmt.Println(utils.Tabulate([]string{"count", "service", "identity", "operation", "url", "code", fmt.Sprintf("traces (limit=%d)", cfg.Limit)}, func() [][]string {
				r := [][]string{}
				for _, i := range results {
					r = append(r, []string{fmt.Sprintf("%d", i.Count), i.Service, i.Identity, i.Operation, i.URL, fmt.Sprintf("%d", i.Code), strings.Join(i.Traces, ",")})
				}
				return r
			}()))

			fmt.Printf("\n> %d results found. You can read the traces from %s/explore and select the jaeger datasource.\n", len(results), cfg.MonitoringURL)
			fmt.Println("  Or run tracer [--stack <name>] --open <trace>.")
		}
	}
}
