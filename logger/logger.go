package logger

import (
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	_ = middleware.RequestLogger(&LogrusFormatter{})
)

func GetLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(NewLogrusFormatter(logger))
}

func NewLogrusFormatter(logger *logrus.Logger) *LogrusFormatter {
	if logger == nil {
		logger = logrus.New()
	}
	return &LogrusFormatter{
		logger: logger,
	}
}

type LogrusFormatter struct {
	logger *logrus.Logger
}

type LogEntry struct {
	flogger logrus.FieldLogger
}

func (l LogEntry) Write(status, bytes int, elapsed time.Duration) {
	switch status / 100 {
	case 2:
		fallthrough
	case 3:
		l.flogger.Info("request completed in ", elapsed)
	case 4:
		l.flogger.Warn("request completed in ", elapsed)
	case 5:
		l.flogger.Error("request completed in ", elapsed)
	}
}

func (l LogEntry) Panic(v interface{}, stack []byte) {
	l.flogger.Error("request failed for: ", string(stack))
}

func guestScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	if v := r.Header.Get("X-FORWARD-PROTO"); v != "" {
		return v
	}
	if v := r.Header.Get("X-FORWARD-SCHEME"); v != "" {
		return v
	}
	return "http"
}

func (l LogrusFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	url := r.URL
	url.Host = r.Host
	if url.Scheme == "" {
		url.Scheme = guestScheme(r)
	}

	return &LogEntry{
		flogger: logrus.NewEntry(l.logger).WithFields(logrus.Fields{
			"URI": url,
		}),
	}
}
