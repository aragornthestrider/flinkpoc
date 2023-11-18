package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	httpReqs = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code, path and method",
		},
		[]string{"code", "path", "method"},
	)
	httpDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_duration_seconds",
			Help:    "Duration of HTTP requests ",
			Buckets: []float64{0.005, 0.05, 0.5, 5, 50},
		},
		[]string{"path", "method"},
	)
)

func (s *Server) encodeJSON(rw http.ResponseWriter, r *http.Request, httpStatus int, v interface{}) {
	rw.Header().Set("content-type", "application/json")
	res, err := json.Marshal(v)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		s.rwFunc(rw, []byte("Internal error"))
	} else {
		rw.WriteHeader(httpStatus)
		s.rwFunc(rw, res)
	}
}

func (s *Server) rwFunc(rw http.ResponseWriter, writeBytes []byte) {
	_, err := rw.Write(writeBytes)
	if err != nil {
		s.Logger.Error("Unable to write", zap.Error(err))
	}
}

type statusWrapWriter struct {
	http.ResponseWriter
	Status int
}

func (r *statusWrapWriter) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusWriter := &statusWrapWriter{ResponseWriter: w, Status: http.StatusOK}
		timer := prometheus.NewTimer(httpDuration.WithLabelValues("/" + strings.Split(r.URL.Path, "/")[0]))
		next.ServeHTTP(statusWriter, r)
		timer.ObserveDuration()
		httpReqs.WithLabelValues(strconv.Itoa(statusWriter.Status), "/"+strings.Split(r.URL.Path, "/")[0])
	})
}
