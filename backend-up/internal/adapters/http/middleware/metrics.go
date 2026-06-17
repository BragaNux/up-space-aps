package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// metricas do Prometheus expostas em /metrics
var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
	ActiveRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_requests",
			Help: "Current number of active HTTP requests",
		},
	)
)

// registra as metricas no Prometheus assim que o pacote carrega
func init() {
	prometheus.MustRegister(RequestsTotal, RequestDuration, ActiveRequests)
}

// WithMetrics envolve o handler pra contar requisicoes, medir duracao e acompanhar requisicoes ativas
func WithMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ActiveRequests.Inc()
		defer ActiveRequests.Dec()

		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)

		status := http.StatusText(rw.status)
		RequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
	})
}

// responseWriter guarda o status code da resposta pra poder usar nas metricas depois
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// MetricsHandler e o handler padrao do Prometheus pra rota /metrics
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
