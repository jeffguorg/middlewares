package sentry

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

const (
	SentryHubCtxKey = "sentry.hub"
)

// Recoverer collects all the panic and report to sentry.
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
				status: 200,
				w:      w,
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
			next.ServeHTTP(fw, r.WithContext(context.WithValue(r.Context(), SentryHubCtxKey, hub)))
		})
	}
}

// GetSentryHub extract sentry hub from http request context
func GetSentryHub(r *http.Request) *sentry.Hub {
	if h := r.Context().Value(SentryHubCtxKey); h != nil {
		return h.(*sentry.Hub)
	}
	return nil
}

func Report(r *http.Request, err error) {
	if h := r.Context().Value(SentryHubCtxKey); h != nil {
		h.(*sentry.Hub).CaptureException(err)
	}
}
