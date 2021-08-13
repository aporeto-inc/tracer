package profiles

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aporeto-inc/tracer/internal/configuration"
	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
)

// Profiles represent a tracer profiles
type Profiles struct {
	Datasources []Datasource `json:"datasources"`
}

// Datasource represent a set of data source
// for a given stack
type Datasource struct {
	LogsIndex                 int    `json:"logsIndex"`
	MetricsIndex              int    `json:"metricsIndex"`
	Name                      string `json:"name"`
	TracesIndex               int    `json:"tracesIndex"`
	TracesDataSourceName      string `json:"tracesDataSourceName"`
	MonitoringCAPath          string `json:"monitoringCAPath"`
	MonitoringCertPath        string `json:"monitoringCertPath"`
	MonitoringCertKeyPath     string `json:"monitoringCertKeyPath"`
	MonitoringCertKeyPassword string `json:"monitoringCertKeyPassword"`
	MonitoringURL             string `json:"monitoringURL"`
}

// PrintDatasources just list the datasources names
func (p Profiles) PrintDatasources() {

	fmt.Println("Here is the list of datasources available:")
	for _, name := range p.Datasources {
		fmt.Printf(" - %s\n", name.Name)
	}

}

// NewProfile will return a new profile using a profile file if exists
// or the arguments if no profile file is set
func NewProfile(cfg *configuration.Configuration) *Datasource {

	path, err := homedir.Expand(cfg.ProfileFile)
	if err != nil {
		zap.L().Fatal("Unable to expand the path", zap.String("path", cfg.ProfileFile), zap.Error(err))
	}

	p, err := parseProfile(path)
	if err != nil {
		zap.L().Fatal("Unable to read profile", zap.String("path", cfg.ProfileFile), zap.Error(err))
	}

	if (p == nil || cfg.MonitoringURL != "") && cfg.Stack == "default" {
		p = &Profiles{}
		p.Datasources = []Datasource{{
			LogsIndex:                 2,
			MetricsIndex:              1,
			Name:                      "default",
			TracesIndex:               3,
			TracesDataSourceName:      "platform-traces",
			MonitoringCAPath:          cfg.MonitoringCAPath,
			MonitoringCertPath:        cfg.MonitoringCertPath,
			MonitoringCertKeyPath:     cfg.MonitoringCertKeyPath,
			MonitoringCertKeyPassword: cfg.MonitoringCertKeyPassword,
			MonitoringURL:             cfg.MonitoringURL,
		}}
	}

	for _, d := range p.Datasources {
		if d.Name == cfg.Stack {

			// Set default if not set
			if d.LogsIndex == 0 {
				if d.MetricsIndex == 0 {
					d.LogsIndex = 2
				} else {
					d.LogsIndex = d.MetricsIndex + 1
				}
			}

			if d.TracesIndex == 0 {
				if d.MetricsIndex == 0 {
					d.TracesIndex = 3
				} else {
					d.TracesIndex = d.MetricsIndex + 2
				}

			}

			if d.MetricsIndex == 0 {
				d.MetricsIndex = 1
			}

			return &d
		}
	}

	zap.L().Error("Unable to find stack name in profile", zap.String("stack", cfg.Stack), zap.String("profile", cfg.ProfileFile))
	p.PrintDatasources()
	os.Exit(1)
	return nil
}

// parseProfile will parse a yaml profile and return a Profile
func parseProfile(profile string) (*Profiles, error) {

	if _, err := os.Stat(profile); err != nil {
		if os.IsNotExist(err) {
			zap.L().Debug("Not profile found", zap.String("path", profile))
			return nil, nil
		}
	}

	zap.L().Debug("Profile found", zap.String("path", profile))

	data, err := ioutil.ReadFile(profile)
	if err != nil {
		return nil, err
	}

	p := &Profiles{}
	if err := yaml.Unmarshal(data, p); err != nil {
		return nil, err
	}

	return p, nil

}
