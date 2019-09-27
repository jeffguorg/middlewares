package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CheckUserCookie(key interface{}, method jwt.SigningMethod) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "user.key", key)
			ctx = context.WithValue(ctx, "user.method", method)

			userCookie, err := r.Cookie("u")
			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			token, err := jwt.Parse(userCookie.Value, func(token *jwt.Token) (interface{}, error) {
				if token.Method != method {
					return nil, fmt.Errorf("Wrong signing method. Expecting %v, got %v", method.Alg(), token.Method.Alg())
				}
				return key, nil
			})
			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx = context.WithValue(r.Context(), "user", token.Claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func MustUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetUser(r) == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUser(r *http.Request) map[string]interface{} {
	v := r.Context().Value("user")
	if v == nil {
		return nil
	}
	return v.(jwt.MapClaims)
}

func UnsetUser(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user",
		MaxAge: -1,
	})
}

func SetUser(w http.ResponseWriter, r *http.Request, user map[string]interface{}) {
	key := r.Context().Value("user.key")
	method := r.Context().Value("user.method")

	if key == nil || method == nil {
		return
	}

	switch method.(type) {
	case jwt.SigningMethod:
		break
	default:
		return
	}

	token := jwt.New(method.(jwt.SigningMethod))
	for k, v := range map[string]interface{}{
		"sub": "backend",
		"iat": float64(time.Now().Unix()),
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	} {
		if _, ok := user[k]; !ok {
			user[k] = v
		}
	}
	token.Claims = jwt.MapClaims(user)
	str, err := token.SignedString(key)
	if err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "user",
		Value: str,
	})
}
