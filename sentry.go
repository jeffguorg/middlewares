package middlewares

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-errors/errors"
	"net/http"
	"os"
)

func SentryLogging(dsn, environment, release string, debug bool) func(handler http.Handler) http.Handler {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         dsn,
		Debug:       debug,
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
					switch rval.(type) {
					case errors.Error:
					case error:
						client.CaptureException(rval.(error), nil, nil)
						break
					default:
						client.CaptureMessage(fmt.Sprint(rval), nil, nil)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
