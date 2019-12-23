package middlewares

import "net/http"

func guessRealIP(r *http.Request) string {
	for _, key := range []string{"X-Real-IP", "X-Forwarded-For"} {
		if val := r.Header.Get(key); len(val) > 0 {
			return val
		}
	}
	return ""
}

func guessScheme(r *http.Request) string {
	for _, key := range []string{"X-Forwarded-Proto", "X-Forwarded-Scheme"} {
		if val := r.Header.Get(key); len(val) > 0 {
			return val
		}
	}
	if len(r.URL.Scheme) > 0 {
		return r.URL.Scheme
	}
	return "http"
}

// ParseProxyHeaders guesses users realip and scheme and such info base on X-* headers
func ParseProxyHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.RemoteAddr) == 0 {
			r.RemoteAddr = guessRealIP(r)
		}
		r.URL.Scheme = guessScheme(r)

		next.ServeHTTP(w, r)
	})
}
