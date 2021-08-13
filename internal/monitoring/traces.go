package monitoring

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/skratchdot/open-golang/open"
	"go.uber.org/zap"
)

// TracingQueryParameters represent the tracing query parameters
type TracingQueryParameters struct {
	End         int64         `url:"end"`
	Limit       int           `url:"limit"`
	Loopback    time.Duration `url:"loopback"`
	MaxDuration time.Duration `url:"maxDuration,omitempty"`
	MinDuration time.Duration `url:"minDuration,omitempty"`
	Service     string        `url:"service"`
	Start       int64         `url:"start"`
	Tags        `url:"tags,omitempty"`
}

// Tags is a tag type with sepcial encoder
type Tags map[string]string

// EncodeValues implement the Encoder interface
func (t Tags) EncodeValues(key string, v *url.Values) error {

	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	v.Add(key, string(data))
	return nil
}

// traceResult is a result of a trace (sparse)
type traceResult struct {
	Data []struct {
		TraceID string `json:"traceID"`
	}
}

// GetTraceIDs try to find traces id related to errors seen in metrics
func (m Client) GetTraceIDs(proxy int, params TracingQueryParameters) ([]string, error) {

	jaegerProxy, err := url.Parse(fmt.Sprintf("api/datasources/proxy/%d/api/traces", proxy))
	if err != nil {
		panic(err)
	}

	q, err := query.Values(params)
	if err != nil {
		return nil, fmt.Errorf("unable to parse query parameters: %w", err)
	}

	request, err := http.NewRequest(http.MethodGet, m.url.ResolveReference(jaegerProxy).String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create new request: %w", err)
	}

	request.URL.RawQuery = q.Encode()

	resp, err := m.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to get traces: %w", err)
	}
	defer resp.Body.Close() // nolint

	p := &traceResult{}
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, fmt.Errorf("unable to decode traces: %w", err)
	}

	zap.L().Debug("Query jaeger", zap.Reflect("params", params), zap.Int("results", len(p.Data)))

	return func() []string {
		res := []string{}
		for _, i := range p.Data {
			res = append(res, i.TraceID)
		}
		return res
	}(), nil
}

// OpenTrace will open a trace in the browser
func OpenTrace(u, datasource, trace string) {

	toOpen, err := url.Parse(u)
	if err != nil {
		fmt.Println("Invalid url: ", err)
		return
	}

	toOpen.Path = "explore"
	q, _ := url.ParseQuery(toOpen.RawQuery)
	q.Add("orgId", "1")
	q.Add("left", fmt.Sprintf(`["now-24h","now","%s",{"query":"%s"}]`, datasource, trace))

	toOpen.RawQuery = q.Encode()

	err = open.Run(toOpen.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Open this URL in your browser:", toOpen)
	}
	os.Exit(0)
}
