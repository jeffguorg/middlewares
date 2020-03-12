package middlewares

import (
	"context"
	"github.com/jeffguorg/middlewares/signature"
	"io/ioutil"
	"net/http"
)

const (
	HttpBodyKey = "httpBodyStore"
)

func StoreBodyInContext() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), HttpBodyKey, body)))
		})
	}
}

func GetBodyContent(request *http.Request) []byte {
	result := request.Context().Value(HttpBodyKey)
	switch resultType := result.(type) {
	case nil:
		return nil
	case string:
		return []byte(resultType)
	case []byte:
		return resultType
	default:
		return nil
	}
}

func CheckSignature(signingMethod signature.SigningMethod, makeSigningString func(r *http.Request) (string, error), getSignature func(r *http.Request) (string, error)) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sign, err := getSignature(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			signingString, err := makeSigningString(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if err := signingMethod.Verify(signingString, sign); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
