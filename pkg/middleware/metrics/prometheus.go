/*
Package prometheus provides middleware to add Prometheus metrics.

Example:
```
package main
import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo-contrib/prometheus"
)
func main() {
    e := echo.New()
    // Enable metrics middleware
    p := prometheus.NewPrometheus("echo", nil)
    p.Use(e)

    e.Logger.Fatal(e.Start(":1323"))
}
```
*/
package metrics

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	goaHttpMdlwr "goa.design/goa/v3/http/middleware"
	goaMdlwr "goa.design/goa/v3/middleware"
)

var defaultMetricPath = "/metrics"
var defaultSubsystem = "goa"

// Standard default metrics
//	counter, counter_vec, gauge, gauge_vec,
//	histogram, histogram_vec, summary, summary_vec
var reqCnt = &Metric{
	ID:          "reqCnt",
	Name:        "requests_total",
	Description: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	Type:        "counter_vec",
	Args:        []string{"code", "method", "host", "url"}}

var reqDur = &Metric{
	ID:          "reqDur",
	Name:        "request_duration_seconds",
	Description: "The HTTP request latency in seconds",
	Type:        "histogram_vec",
	Args:        []string{"method", "host", "url"},
}

var resSz = &Metric{
	ID:          "resSz",
	Name:        "response_size_bytes",
	Description: "The HTTP response sizes in bytes.",
	Type:        "summary"}

var reqSz = &Metric{
	ID:          "reqSz",
	Name:        "request_size_bytes",
	Description: "The HTTP request sizes in bytes.",
	Type:        "summary"}

var standardMetrics = []*Metric{
	reqCnt,
	reqDur,
	resSz,
	reqSz,
}

/*
RequestCounterURLLabelMappingFunc is a function which can be supplied to the middleware to control
the cardinality of the request counter's "url" label, which might be required in some contexts.
For instance, if for a "/customer/:name" route you don't want to generate a time series for every
possible customer name, you could use this function:

func(c echo.Context) string {
	url := c.Request.URL.Path
	for _, p := range c.Params {
		if p.Key == "name" {
			url = strings.Replace(url, p.Value, ":name", 1)
			break
		}
	}
	return url
}

which would map "/customer/alice" and "/customer/bob" to their template "/customer/:name".
*/
type RequestCounterURLLabelMappingFunc func(c echo.Context) string

// Metric is a definition for the name, description, type, ID, and
// prometheus.Collector type (i.e. CounterVec, Summary, etc) of each metric
type Metric struct {
	MetricCollector prometheus.Collector
	ID              string
	Name            string
	Description     string
	Type            string
	Args            []string
}

// Prometheus contains the metrics gathered by the instance and its path
type Prometheus struct {
	reqCnt        *prometheus.CounterVec
	reqSz, resSz  prometheus.Summary
	reqDur        *prometheus.HistogramVec
	router        *echo.Echo
	listenAddress string
	Ppg           PushGateway

	MetricsList []*Metric
	MetricsPath string
	Subsystem   string
	Skipper     middleware.Skipper

	RequestCounterURLLabelMappingFunc RequestCounterURLLabelMappingFunc

	// Context string to use as a prometheus URL label
	URLLabelFromContext string
}

// PushGateway contains the configuration for pushing to a Prometheus pushgateway (optional)
type PushGateway struct {
	// Push interval in seconds
	PushIntervalSeconds time.Duration

	// Push Gateway URL in format http://domain:port
	// where JOBNAME can be any string of your choice
	PushGatewayURL string

	// Local metrics URL where metrics are fetched from, this could be ommited in the future
	// if implemented using prometheus common/expfmt instead
	MetricsURL string

	// pushgateway job name, defaults to "echo"
	Job string
}

// NewPrometheus generates a new set of metrics with a certain subsystem name
func NewPrometheus(subsystem string, skipper middleware.Skipper, customMetricsList ...[]*Metric) *Prometheus {
	var metricsList []*Metric
	if skipper == nil {
		skipper = middleware.DefaultSkipper
	}

	if len(customMetricsList) > 1 {
		panic("Too many args. NewPrometheus( string, <optional []*Metric> ).")
	} else if len(customMetricsList) == 1 {
		metricsList = customMetricsList[0]
	}

	metricsList = append(metricsList, standardMetrics...)

	p := &Prometheus{
		MetricsList: metricsList,
		MetricsPath: defaultMetricPath,
		Subsystem:   defaultSubsystem,
		Skipper:     skipper,
		RequestCounterURLLabelMappingFunc: func(c echo.Context) string {
			return c.Path() // i.e. by default do nothing, i.e. return URL as is
		},
	}

	p.registerMetrics(subsystem)

	return p
}

// SetPushGateway sends metrics to a remote pushgateway exposed on pushGatewayURL
// every pushIntervalSeconds. Metrics are fetched from metricsURL
func (p *Prometheus) SetPushGateway(pushGatewayURL, metricsURL string, pushIntervalSeconds time.Duration) {
	p.Ppg.PushGatewayURL = pushGatewayURL
	p.Ppg.MetricsURL = metricsURL
	p.Ppg.PushIntervalSeconds = pushIntervalSeconds
	p.startPushTicker()
}

// SetPushGatewayJob job name, defaults to "echo"
func (p *Prometheus) SetPushGatewayJob(j string) {
	p.Ppg.Job = j
}

// SetListenAddress for exposing metrics on address. If not set, it will be exposed at the
// same address of the echo engine that is being used
// func (p *Prometheus) SetListenAddress(address string) {
// 	p.listenAddress = address
// 	if p.listenAddress != "" {
// 		p.router = echo.Echo().Router()
// 	}
// }

// SetListenAddressWithRouter for using a separate router to expose metrics. (this keeps things like GET /metrics out of
// your content's access log).
// func (p *Prometheus) SetListenAddressWithRouter(listenAddress string, r *echo.Echo) {
// 	p.listenAddress = listenAddress
// 	if len(p.listenAddress) > 0 {
// 		p.router = r
// 	}
// }

func (p *Prometheus) Start(address string) {
	e := echo.New()
	p.SetMetricsPath(e)
	e.Logger.Fatal(e.Start(address))
}

// SetMetricsPath set metrics paths
func (p *Prometheus) SetMetricsPath(e *echo.Echo) {
	if p.listenAddress != "" {
		p.router.GET(p.MetricsPath, prometheusHandler())
		p.runServer()
	} else {
		e.GET(p.MetricsPath, prometheusHandler())
	}
}

func (p *Prometheus) runServer() {
	if p.listenAddress != "" {
		go func(p *Prometheus) {
			_ = p.router.Start(p.listenAddress)
		}(p)
	}
}

func (p *Prometheus) getMetrics() []byte {
	response, err := http.Get(p.Ppg.MetricsURL)
	if err != nil {
		log.Errorf("Error getting metrics: %v", err)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return body
}

func (p *Prometheus) getPushGatewayURL() string {
	h, _ := os.Hostname()
	if p.Ppg.Job == "" {
		p.Ppg.Job = "echo"
	}
	return p.Ppg.PushGatewayURL + "/metrics/job/" + p.Ppg.Job + "/instance/" + h
}

func (p *Prometheus) sendMetricsToPushGateway(metrics []byte) {
	req, _ := http.NewRequest("POST", p.getPushGatewayURL(), bytes.NewBuffer(metrics))
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error sending to push gateway: %v", err)
	}

	resp.Body.Close()
}

func (p *Prometheus) startPushTicker() {
	ticker := time.NewTicker(time.Second * p.Ppg.PushIntervalSeconds)
	go func() {
		for range ticker.C {
			p.sendMetricsToPushGateway(p.getMetrics())
		}
	}()
}

// NewMetric associates prometheus.Collector based on Metric.Type
func NewMetric(m *Metric, subsystem string) prometheus.Collector {
	var metric prometheus.Collector
	switch m.Type {
	case "counter_vec":
		metric = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "counter":
		metric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "gauge_vec":
		metric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "gauge":
		metric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "histogram_vec":
		metric = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "histogram":
		metric = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "summary_vec":
		metric = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "summary":
		metric = prometheus.NewSummary(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	}
	return metric
}

func (p *Prometheus) registerMetrics(subsystem string) {
	for _, metricDef := range p.MetricsList {
		metric := NewMetric(metricDef, subsystem)
		if err := prometheus.Register(metric); err != nil {
			log.Errorf("%s could not be registered in Prometheus: %v", metricDef.Name, err)
		}
		switch metricDef {
		case reqCnt:
			p.reqCnt = metric.(*prometheus.CounterVec)
		case reqDur:
			// p.reqDur = metric.(prometheus.Summary)
			p.reqDur = metric.(*prometheus.HistogramVec)
		case resSz:
			p.resSz = metric.(prometheus.Summary)
		case reqSz:
			p.reqSz = metric.(prometheus.Summary)
		}
		metricDef.MetricCollector = metric
	}
}

func (p *Prometheus) HandlerFunc(l goaMdlwr.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Context().Value(goaMdlwr.RequestIDKey)
			if reqID == nil {
				reqID = shortID()
			}

			started := time.Now()
			reqSz := computeApproximateRequestSize(r)

			_ = l.Log("id", reqID,
				"req", r.Method+" "+r.URL.String(),
				"from", from(r))

			rw := goaHttpMdlwr.CaptureResponse(w)
			h.ServeHTTP(rw, r)

			_ = l.Log("id", reqID,
				"req", r.Method+" "+r.URL.String(),
				"from", from(r),
				"status", rw.StatusCode,
				"bytes", rw.ContentLength,
				"time", time.Since(started).String())

			elapsed := float64(time.Since(started)) / float64(time.Second)

			// p.reqDur.Observe(elapsed)
			p.reqDur.WithLabelValues(r.Method, r.Host, r.URL.EscapedPath()).Observe(elapsed)

			status := strconv.Itoa(rw.StatusCode)

			p.reqCnt.WithLabelValues(status, r.Method, r.Host, r.URL.EscapedPath()).Inc()
			p.reqSz.Observe(float64(reqSz))
			p.resSz.Observe(float64(rw.ContentLength))
		})
	}
}

func prometheusHandler() echo.HandlerFunc {
	h := promhttp.Handler()
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)

	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}

// shortID produces a " unique" 6 bytes long string.
// Do not use as a reliable way to get unique IDs, instead use for things like logging.
func shortID() string {
	b := make([]byte, 6)
	_, _ = io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// from makes a best effort to compute the request client IP.
func from(req *http.Request) string {
	if f := req.Header.Get("X-Forwarded-For"); f != "" {
		return f
	}
	f := req.RemoteAddr
	ip, _, err := net.SplitHostPort(f)
	if err != nil {
		return f
	}
	return ip
}
