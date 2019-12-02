package middlewares

import (
	"context"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	ctxPrefix = "IsylLzqZ"
	json      = jsoniter.ConfigCompatibleWithStandardLibrary
)

// RequireParametersInQuery checks for parameters existence in query string
func RequireParametersInQuery(keys ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			for _, key := range keys {
				if len(r.URL.Query().Get(key)) == 0 {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
			for k := range r.URL.Query() {
				ctx = context.WithValue(ctx, ctxPrefix+k, r.URL.Query().Get(k))
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireParametersInQuery checks for parameters existence in query string
func RequireParametersInJSON(keys ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var values map[string]interface{}
			var ctx = r.Context()

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
			for k, v := range values {
				ctx = context.WithValue(ctx, ctxPrefix+k, v)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Parameter(r *http.Request, k string) interface{} {
	return r.Context().Value(ctxPrefix + k)
}

func ParameterStringWithDefault(r *http.Request, k string, d string) string {
	if param := r.Context().Value(ctxPrefix + k); param != nil {
		if v, ok := param.(string); ok {
			return v
		}
		return d
	}
	return d
}

func ParameterIntWithDefault(r *http.Request, k string, d int) int {
	if param := r.Context().Value(ctxPrefix + k); param != nil {
		switch v := param.(type) {
		case int:
			return int(v)
		case int8:
			return int(v)
		case int16:
			return int(v)
		case int32:
			return int(v)
		case int64:
			return int(v)
		case uint:
			return int(v)
		case uint8:
			return int(v)
		case uint16:
			return int(v)
		case uint32:
			return int(v)
		case uint64:
			return int(v)
		case string:
			if n, err := strconv.Atoi(v); err == nil {
				return n
			}
		}
	}
	return d
}
