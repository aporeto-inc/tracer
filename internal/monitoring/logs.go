// Package monitoring provides loki related monitoring logic
package monitoring

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aporeto-inc/tracer/internal/configuration"
	"github.com/grafana/loki/pkg/logcli/client"
	"github.com/grafana/loki/pkg/logcli/output"
	"github.com/grafana/loki/pkg/logcli/query"
	"github.com/prometheus/common/config"
)

// GetLogs try to get the logs for a service and and a time window
func (m Client) GetLogs(proxy int, from, to time.Time, services []string, cfg configuration.LogConf, quiet bool) error {
	lokiProxy, err := m.url.Parse(fmt.Sprintf("api/datasources/proxy/%d", proxy))
	if err != nil {
		panic(err)
	}

	client := &client.DefaultClient{
		TLSConfig: config.TLSConfig{
			CAFile:   m.cfg.MonitoringCAPath,
			CertFile: m.cfg.MonitoringCertPath,
			KeyFile:  m.cfg.MonitoringCertKeyPath,
		},
		Address: m.url.ResolveReference(lokiProxy).String(),
	}

	q := &query.Query{
		QueryString: func() string {
			if len(services) > 0 {
				return fmt.Sprintf(`{app=~"%s"} %s`, strings.Join(services, "|"), cfg.LogFilter)
			}
			return cfg.LogFilter
		}(),
		Start:           from,
		End:             to,
		Limit:           cfg.LogLines,
		BatchSize:       100,
		Quiet:           quiet,
		NoLabels:        cfg.LogNoLabels,
		ShowLabelsKey:   []string{"pod"},
		IgnoreLabelsKey: []string{"filename", "stream"},
		ColoredOutput:   true,
	}

	outputOptions := &output.LogOutputOptions{
		NoLabels:      q.NoLabels,
		ColoredOutput: q.ColoredOutput,
	}

	if cfg.Direction == "forward" {
		q.Forward = true
	}

	out, err := output.NewLogOutput(os.Stdout, "default", outputOptions)
	if err != nil {
		return err
	}

	if cfg.Follow {
		q.TailQuery(0, client, out)
	} else {
		q.DoQuery(client, out, false)
	}

	return nil
}
