package middlewares

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// RequireParametersInQuery checks for parameters existence in query string
func RequireParametersInQuery(keys ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, key := range keys {
				if len(r.URL.Query().Get(key)) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireParametersInQuery checks for parameters existence in query string
func RequireParametersInJSON(keys ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var values map[string]interface{}

			buf, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err := json.Unmarshal(buf, &values); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			for _, key := range keys {
				if _, ok := values[key]; !ok {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

			}
			next.ServeHTTP(w, r)
		})
	}
}
