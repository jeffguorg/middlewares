package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/go-errors/errors"
	"net/http"
)

var (
	ErrCSRF = errors.New("CSRF Detected")
)

func digest(key, content []byte) string {
	hasher := hmac.New(sha256.New, key)
	hasher.Write(content)
	digest := hasher.Sum(nil)
	return base64.StdEncoding.EncodeToString(digest)
}

func SetSecretCookie(w http.ResponseWriter, cookiename string, content, key []byte) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookiename,
		Value:  digest(key, content),
		Secure: true,
	})
}

// RequireParametersInQuery checks for parameters existence in query string
func SecureCookie(key []byte, queryname, cookiename string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			orig := r.URL.Query().Get(queryname)
			cookie, err := r.Cookie(cookiename)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			sign := digest(key, []byte(orig))
			if cookie.Value != sign {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func CheckSecureCookie(csrf string, orig, key []byte) error {
	sign := digest(key, []byte(orig))
	if sign == csrf {
		return nil
	}
	return nil
}
