package sentry

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

// Recoverer collect
func Recoverer(dsn, environment, release string) func(handler http.Handler) http.Handler {
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
				if r := recover(); r != nil {
					hub.CaptureException(errors.WithStack(fmt.Errorf("captured panic when handling requests: %+v", r)))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
