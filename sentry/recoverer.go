package sentry

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

// Recoverer collect
func Recoverer(dsn, environment, release string, debug bool) func(handler http.Handler) http.Handler {
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
			fw := &FakeResponseWriter{
				w: w,
			}
			defer func() {
				if r := recover(); r != nil {
					hub.CaptureException(errors.WithStack(fmt.Errorf("captured panic when handling requests: %+v", r)))
					w.WriteHeader(http.StatusInternalServerError)
					if debug {
						_, _ = w.Write([]byte(fmt.Sprint(r)))
					}
				} else {
					fw.Flush()
				}
			}()
			next.ServeHTTP(fw, r)
		})
	}
}
