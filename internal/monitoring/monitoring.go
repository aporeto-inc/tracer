package monitoring

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aporeto-inc/tracer/internal/profiles"
	"go.aporeto.io/tg/tglib"
)

// Client is a monitoring client that can query the monitoring stacks
type Client struct {
	url    *url.URL
	cfg    *profiles.Datasource
	client http.Client
}

// NewClient return a new montitoring.Client
func NewClient(cfg *profiles.Datasource) (*Client, error) {
	url, err := url.Parse(cfg.MonitoringURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url: %w", err)
	}

	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("cannot create system cert pool: %w", err)
	}

	if cfg.MonitoringCAPath != "" {

		publicCACertData, err := os.ReadFile(cfg.MonitoringCAPath)
		if err != nil {
			return nil, fmt.Errorf("cannot read provided ca: %w", err)
		}

		// In case the pool is empty
		if pool == nil {
			pool = x509.NewCertPool()
		}

		if ok := pool.AppendCertsFromPEM(publicCACertData); !ok {
			return nil, fmt.Errorf("cannot append provided ca certificate: %w", err)
		}
	}

	x509ClientCert, pkey, err := tglib.ReadCertificatePEM(cfg.MonitoringCertPath, cfg.MonitoringCertKeyPath, cfg.MonitoringCertKeyPassword)
	if err != nil {
		return nil, fmt.Errorf("cannot read provided monitoring certificate: %w", err)
	}
	clientCert, err := tglib.ToTLSCertificate(x509ClientCert, pkey)
	if err != nil {
		return nil, fmt.Errorf("cannot convert provided monitoring certificate: %w", err)
	}

	return &Client{
		client: http.Client{
			Timeout: 120 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      pool,
					Certificates: []tls.Certificate{clientCert},
				},
			},
		},
		url: url, cfg: cfg,
	}, nil
}
