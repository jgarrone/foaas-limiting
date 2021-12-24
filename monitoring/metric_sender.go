package monitoring

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type LatencyMetricSender interface {
	Send(start, end time.Time, method, status string)
}

type httpLatencyMetricSender struct {
}

func NewHttpLatencyMetricSender() *httpLatencyMetricSender {
	return &httpLatencyMetricSender{}
}

// Send currently logs the time it took to complete some method and its status. A more interesting
// implementation would be to send metrics to an observability service.
func (h httpLatencyMetricSender) Send(start, end time.Time, method, status string) {
	elapsedMs := end.Sub(start).Milliseconds()
	log.Debugf("Method %q took %dms to complete, status: %s", method, elapsedMs, status)
}
