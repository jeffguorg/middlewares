package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

func SetSecretCookie(w http.ResponseWriter, cookiename string, content, key []byte) {
	hasher := hmac.New(sha256.New, key)
	hasher.Write([]byte(content))
	digest := hasher.Sum(nil)
	sign := base64.StdEncoding.EncodeToString(digest)
	http.SetCookie(w, &http.Cookie{
		Name:   cookiename,
		Value:  sign,
		Secure: true,
	})
}

// RequireParametersInQuery checks for parameters existence in query string
func CheckSecretCookie(key []byte, queryname, cookiename string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			orig := r.URL.Query().Get(queryname)
			cookie, err := r.Cookie(cookiename)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			hasher := hmac.New(sha256.New, key)
			hasher.Write([]byte(orig))
			digest := hasher.Sum(nil)
			sign := base64.StdEncoding.EncodeToString(digest)

			if cookie.Value != sign {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
