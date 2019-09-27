package middlewares

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
	return middleware.RequestLogger(NewLogrusFormattter(logger))
}

func NewLogrusFormattter(logger *logrus.Logger) *LogrusFormatter {
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
	start   time.Time
}

func (l LogEntry) Write(status, bytes int, elapsed time.Duration) {
	l.flogger.Info("request completed in ", time.Now().Sub(l.start).Seconds(), " seconds")
}

func (l LogEntry) Panic(v interface{}, stack []byte) {
	l.flogger.Error("request failed for: ", string(stack))
}

func (l LogrusFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &LogEntry{
		flogger: logrus.NewEntry(l.logger).WithFields(logrus.Fields{}),
		start:   time.Now(),
	}
}
