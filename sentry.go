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
	cli, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         dsn,
		ServerName:  hostname,
		Release:     release,
		Environment: environment,
	})
	if err != nil {
		panic(err)
	}
	hub := sentry.NewHub(cli, sentry.NewScope())
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rval := recover(); rval != nil {
					hub.CaptureException(errors.New(rval))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
