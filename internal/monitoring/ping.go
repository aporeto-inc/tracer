package monitoring

import (
	"fmt"
	"net/http"
	"net/url"
)

// Ping check if the monitoring endpoint is reachable
func (m *Client) Ping() error {

	url, err := url.Parse("/api/health")
	if err != nil {
		panic(err)
	}

	req, err := m.client.Get(m.url.ResolveReference(url).String())
	if err != nil {
		return fmt.Errorf("unable to query monitoring url: %w", err)
	}
	defer req.Body.Close() //nolint
	if req.StatusCode != http.StatusOK {
		return fmt.Errorf("monitoring doesnt seems helhty: return code %d", req.StatusCode)
	}

	return nil
}
