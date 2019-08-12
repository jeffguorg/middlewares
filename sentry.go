package middlewares

import (
	"github.com/getsentry/sentry-go"
	"github.com/go-errors/errors"
	"net/http"
	"os"
)

func SentryLogging(dsn, environment, release string) func(handler http.Handler) http.Handler {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	err = sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		ServerName:  hostname,
		Release:     release,
		Environment: environment,
	})
	if err != nil {
		panic(err)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rval := recover(); rval != nil {
					sentry.CaptureException(errors.New(rval))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
