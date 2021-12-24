package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jgarrone/foaas-limiting/monitoring"
)

const ErrorLabel = "error"

// Client implements a wrapper for an http client that automatically sends latency metrics.
type Client struct {
	httpClient          Doer
	latencyMetricSender monitoring.LatencyMetricSender
}

func NewClient(httpClient Doer, latencyMetricSender monitoring.LatencyMetricSender) *Client {
	return &Client{
		httpClient:          httpClient,
		latencyMetricSender: latencyMetricSender,
	}
}

func (c *Client) Do(r *http.Request, methodName string) (*http.Response, error) {
	start := time.Now()
	resp, err := c.httpClient.Do(r)
	end := time.Now()

	c.latencyMetricSender.Send(start, end, methodName, c.statusLabelFrom(resp))

	return resp, err
}

func (*Client) statusLabelFrom(resp *http.Response) string {
	if resp == nil {
		return ErrorLabel
	}

	return strconv.Itoa(resp.StatusCode)
}
