package session

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	defaultSessionCtxKey       = "session.ctx.session"
	defaultSessionChangeCtxKey = "session.ctx.change"
	defaultSessionDuration     = time.Minute * 30

	// InfiniteSessionDuration means that the session won't expire unless you close the browser
	InfiniteSessionDuration time.Duration = 0
)

// Client save or load session to a storage
type Client interface {
	Update(sessionID string, session map[string]interface{}) error
	Reset(sessionID string) error
	Load(sessionID string) (map[string]interface{}, error)
}

// Mixin provides session context and functionality
type Mixin struct {
	client              Client
	sessionCtxKey       interface{}
	sessionChangeCtxKey interface{}

	CookieName string
	CookiePath string

	SessionDuration   time.Duration
	JWTKey            string
	JWTSessionKeyname string
}

// MixinOption sets the option in mixin
type MixinOption func(*Mixin)

// SetMixinCtxKey configure context key
func SetMixinCtxKey(sessionKey, sessionChangeKey interface{}) MixinOption {
	return func(m *Mixin) {
		m.sessionCtxKey = sessionKey
		m.sessionChangeCtxKey = sessionChangeKey
	}
}

// SetDuration configure session's duration
func SetDuration(duration time.Duration) MixinOption {
	return func(m *Mixin) {
		m.SessionDuration = duration
	}
}

// NewMixin return a configured mixin to use as middleware
func NewMixin(client Client, cookieName, cookiePath, jwtKey, jwtSessionKeyName string, options ...MixinOption) Mixin {
	mixin := Mixin{
		CookieName: cookieName,
		CookiePath: cookiePath,

		JWTKey:            jwtKey,
		JWTSessionKeyname: jwtSessionKeyName,

		client:              client,
		sessionCtxKey:       defaultSessionCtxKey,
		sessionChangeCtxKey: defaultSessionChangeCtxKey,
		SessionDuration:     defaultSessionDuration,
	}

	for _, opt := range options {
		opt(&mixin)
	}

	return mixin
}

// EnsureSession ensures a session storage in context
func (mixin Mixin) EnsureSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var sessionID string
		var session map[string]interface{}
		var sessionChange map[string]interface{} = make(map[string]interface{})
		var jwtToken *jwt.Token

		sessionCookie, err := r.Cookie(mixin.CookieName)
		// if there is already a session storage
		if err == nil {
			// try parse it
			jwtToken, err = jwt.ParseWithClaims(sessionCookie.Value, make(jwt.MapClaims), func(token *jwt.Token) (interface{}, error) {
				if token.Header["alg"] != jwt.SigningMethodHS256.Alg() {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(mixin.JWTKey), nil
			})
			// if there is no issue with the token

			if err == nil {
				if claimSessionID, ok := (jwtToken.Claims.(jwt.MapClaims))[mixin.JWTSessionKeyname]; ok {
					session, err = mixin.client.Load(claimSessionID.(string))
				}
			}
		}

		if err != nil {
			// no session is found
			sessionID, _ := uuid.NewRandom()
			claims := jwt.MapClaims{
				"iat":                   time.Now(),
				mixin.JWTSessionKeyname: sessionID.String(),
			}
			if mixin.SessionDuration != 0 {
				claims["exp"] = time.Now().Add(mixin.SessionDuration)
			}
			jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			tokenStr, err := jwtToken.SignedString([]byte(mixin.JWTKey))
			if err != nil {
				panic(err)
			}

			// set session id to cookie
			sessionCookie = &http.Cookie{
				Name:  mixin.CookieName,
				Value: tokenStr,
				Path:  mixin.CookiePath,
			}
			if mixin.SessionDuration != 0 {
				sessionCookie.Expires = time.Now().Add(mixin.SessionDuration)
			}
			http.SetCookie(rw, sessionCookie)

			// init empty session
			session = make(map[string]interface{})
		}

		ctx = context.WithValue(ctx, mixin.sessionCtxKey, session)
		ctx = context.WithValue(ctx, mixin.sessionChangeCtxKey, sessionChange)

		defer mixin.client.Update(sessionID, sessionChange)

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

// Set update the session when the request end
func (mixin Mixin) Set(request *http.Request, name string, value interface{}) {
	ctxValue := request.Context().Value(mixin.sessionChangeCtxKey)
	if sessionChange, ok := ctxValue.(map[string]interface{}); ok {
		sessionChange[name] = value
	}
}

// Get returns the value in session
func (mixin Mixin) Get(request *http.Request, name string) interface{} {
	ctxValue := request.Context().Value(mixin.sessionCtxKey)
	if session, ok := ctxValue.(map[string]interface{}); ok {
		if value, ok := session[name]; ok {
			return value
		}
	}
	return nil
}
