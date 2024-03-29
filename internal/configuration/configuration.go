// Package configuration is a small package for handling configuration
package configuration

import (
	"fmt"
	"time"

	"go.aporeto.io/addedeffect/lombric"
	"go.aporeto.io/underwater/logutils"
)

var (
	version = "v0.0.0"
	commit  = "dev"
)

// MonitoringConf hold the grafana configuration
// This is used to proxy all requests to datasources
type MonitoringConf struct {
	MonitoringCAPath          string `mapstructure:"monitoring-ca-path"            desc:"Path to the monitoring CA certificate" `
	MonitoringURL             string `mapstructure:"monitoring-url"                desc:"The monitoring url to use"             `
	MonitoringCertPath        string `mapstructure:"monitoring-cert"               desc:"Path to the monitoring cert"           `
	MonitoringCertKeyPath     string `mapstructure:"monitoring-cert-key"           desc:"Path to the monitoring cert key"       `
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

// LogConf is the configuration realted to logs
type LogConf struct {
	Direction   string `mapstructure:"direction" desc:"Logs: Direction of the logs" default:"forward" allowed:"forward,backward"`
	LogFilter   string `mapstructure:"log-filter" desc:"Logs; Optional log filter to append to log query if service flag is used or full logcli filter if not service flag are set"`
	LogLines    int    `mapstructure:"lines" desc:"Logs: Number of lines to print" default:"10"`
	Log         bool   `mapstructure:"log" desc:"Logs: Enable log mode to get logs from services"`
	Follow      bool   `mapstructure:"follow" desc:"Logs: Follow logs stream in almost real time"`
	LogNoLabels bool   `mapstructure:"no-labels" desc:"Logs: Do not display labels with logs"`
}

// Configuration hold the service configuration.
type Configuration struct {
	MonitoringConf `mapstructure:",squash"`
	LoggingConf    `mapstructure:",squash"`
	Open           string `mapstructure:"open" desc:"Traces: Open a given trace to your browser."`
	ProfileFile    string `mapstructure:"profile-file" desc:"Profile file: the profile file pathto use." default:"~/.tracer/default.yaml"`
	Stack          string `mapstructure:"stack" desc:"Stack: The stack name to use if any." default:"default"`
	FilterConf     `mapstructure:",squash"`
	TimeWindow     `mapstructure:",squash"`
	LogConf        `mapstructure:",squash"`
	TraceConf      `mapstructure:",squash"`
	Help           bool `mapstructure:"help" desc:"Show full help with examples"`
}

// Prefix returns the configuration prefix.
func (c *Configuration) Prefix() string { return "tracer" }

// PrintVersion prints the current version.
func (c *Configuration) PrintVersion() {
	fmt.Printf("tracer - %s (%s)\n", version, commit)
}

// NewConfiguration returns a new configuration.
func NewConfiguration() *Configuration {
	c := &Configuration{}
	lombric.Initialize(c)
	logutils.Configure(c.LogLevel, c.LogFormat)

	if c.Help {
		showHelp()
	}

	return c
}
