package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

var (
	_ = middleware.RequestLogger(&LogrusFormatter{})
)

// GetLogger make a http compatible middleware that logs every requests
func GetLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	formatter := NewLogrusFormatter(logger)
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := formatter.NewLogEntry(r)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				if r := recover(); r != nil {
					entry.Panic(r, nil)
				} else {
					status := ww.Status()
					if status == 0 {
						status = 200
					}
					entry.Write(status, ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
				}
			}()

			next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))
		}
		return http.HandlerFunc(fn)
	}
}

// NewLogrusFormatter return a formatter that write log to specified destination
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

func (l LogEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	l.flogger = l.flogger.WithFields(logrus.Fields{
		"status":  status,
		"written": bytes,
	})
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

func (l LogEntry) Panic(panic interface{}, stack []byte) {
	l.flogger.Errorf("panic occurred for '%v' at %v", panic, string(stack))
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
			"uri":    url,
			"method": r.Method,
		}),
	}
}
