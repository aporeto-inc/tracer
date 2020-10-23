// Package configuration is a small package for handling configuration
package configuration

import (
	"fmt"
	"time"

	"go.aporeto.io/addedeffect/lombric"
)

// MonitoringConf hold the grafana configuration
// This is used to proxy all requests to datasources
type MonitoringConf struct {
	MonitoringCAPath          string `mapstructure:"monitoring-ca-path"            desc:"Path to the monitoring CA certificate" `
	MonitoringURL             string `mapstructure:"monitoring-url"                desc:"The monitoring url to use"             required:"true"`
	MonitoringCertPath        string `mapstructure:"monitoring-cert"               desc:"Path to the monitoring cert"           required:"true"`
	MonitoringCertKeyPath     string `mapstructure:"monitoring-cert-key"           desc:"Path to the monitoring cert key"       required:"true"`
	MonitoringCertKeyPassword string `mapstructure:"monitoring-cert-key-pass"      desc:"Password for the monitoring cert key"  secret:"true"`
}

// LoggingConf is the configuration for log.
type LoggingConf struct {
	LogFormat string `mapstructure:"log-format"   desc:"Log format"   default:"console"`
	LogLevel  string `mapstructure:"log-level"    desc:"Log level"    default:"info"`
}

// TimeWindow is the configuration for queries
type TimeWindow struct {
	From  string        `mapstructure:"from"   desc:"From date"`
	To    string        `mapstructure:"to"     desc:"To date"`
	Since time.Duration `mapstructure:"since"  desc:"Since duration (will compute From and To with currrent date)" default:"1h"`
}

// TraceConf is the configuration related to traces
type TraceConf struct {
	Namespace   string        `mapstructure:"namespace" desc:"Traces: Lookg for traces matching that namespace"`
	OnlyError   bool          `mapstructure:"errors-only" desc:"Traces: Look only for trace in error"`
	MinDuration time.Duration `mapstructure:"slower-than" desc:"Traces: Look for traces slower than the provided duration"`
	Limit       int           `mapstructure:"limit" desc:"Traces: The number of traces to display" default:"1"`
}

// FilterConf is the configuration realted to filters
type FilterConf struct {
	Codes    string   `mapstructure:"code" desc:"Filters: The code to filter ex:200-300,400-422,500"`
	Services []string `mapstructure:"service" desc:"Filters: The service to filter (repeatable)"`
	URLS     []string `mapstructure:"url" desc:"Filters: The url to filter (repeatable)"`
}

// Configuration hold the service configuration.
type Configuration struct {
	LoggingConf    `mapstructure:",squash"`
	FilterConf     `mapstructure:",squash"`
	MonitoringConf `mapstructure:",squash"`
	TimeWindow     `mapstructure:",squash"`
	TraceConf      `mapstructure:",squash"`

	Help bool   `mapstructure:"help" desc:"Show full help with examples"`
	Open string `mapstructure:"open" desc:"Open a given trace to your browser."`
}

// Prefix returns the configuration prefix.
func (c *Configuration) Prefix() string { return "tracer" }

// PrintVersion prints the current version.
func (c *Configuration) PrintVersion() {
	fmt.Printf("tracer - %s (%s)\n", ProjectVersion, ProjectSha)
}

// NewConfiguration returns a new configuration.
func NewConfiguration() *Configuration {

	c := &Configuration{}
	lombric.Initialize(c)

	if c.Help {
		showHelp()
	}

	if c.Open != "" {
		openTrace(c.MonitoringURL, c.Open)
	}

	return c
}
