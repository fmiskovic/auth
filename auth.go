package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Handle returns an HTTP middleware that validates JWT tokens from the Authorization header.
// It uses the provided configuration and options to create a new Auth instance and returns its Handler.
//
// Example usage:
//
//	cfg := auth.Config{
//	    Secret:   "your-secret",
//	    TokenExp: time.Hour,
//	}
//
//	mux := http.NewServeMux()
//
//	handler := func(w http.ResponseWriter, r *http.Request) {
//	    w.WriteHeader(http.StatusOK)
//	    w.Write([]byte("success"))
//	}
//
//	mux.Handle("/protected", auth.Handle(cfg)(http.HandlerFunc(handler)))
//
//	http.ListenAndServe(":8080", mux)
//
// Parameters:
//
//	cfg - the configuration for the Auth instance.
//	opts - optional parameters for customizing the Auth instance.
//
// Returns:
//
//	func(next http.Handler) http.Handler - the middleware function that wraps the next http.Handler with JWT validation.
func Handle(cfg Config, opts ...Option) func(next http.Handler) http.Handler {
	a := New(cfg, opts...)
	return a.Handler
}

type Auth struct {
	cfg  Config
	opts *options
}

func New(cfg Config, opts ...Option) *Auth {
	o := new(options)
	for _, opt := range opts {
		opt(o)
	}

	return &Auth{cfg: cfg, opts: o}
}

// Handler is an HTTP middleware that validates JWT tokens from the Authorization header.
// If the token is valid, it extracts the claims and adds them to the request context.
// If the token is invalid or missing, it responds with an unauthorized status.
// If the admin option is enabled, it also checks if the user has admin privileges.
//
// Example usage:
//
//	cfg := auth.Config{
//	    Secret:   "your-secret",
//	    TokenExp: time.Hour,
//	}
//
//	mux := http.NewServeMux()
//
//	handler := func(w http.ResponseWriter, r *http.Request) {
//	    w.WriteHeader(http.StatusOK)
//	    w.Write([]byte("success"))
//	}
//
//	mux.Handle("/protected", auth.New(cfg).Handler(http.HandlerFunc(handler)))
//
//	http.ListenAndServe(":8080", mux)
//
// Parameters:
//
//	next - the next http.Handler to be called if the token is valid.
//
// Returns:
//
//	http.Handler - the wrapped http.Handler with JWT validation.
func (a *Auth) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := a.extractBearerToken(r)
		if token == "" {
			http.Error(w, "unauthorized request", http.StatusUnauthorized)
			return
		}

		claims, err := a.validateToken(token)
		if err != nil {
			http.Error(
				w,
				fmt.Errorf("unauthorized request: %v", err).Error(),
				http.StatusUnauthorized,
			)
		}

		ctx := claims.toContext()
		if a.opts.admin && !ctx.IsAdmin { // check if the user is admin
			http.Error(w, "forbidden request", http.StatusForbidden)
			return
		}

		// add claims to the request context and pass it to the next handler
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKey, ctx)))
	})
}

func (a *Auth) extractBearerToken(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

func (a *Auth) validateToken(token string) (*JwtClaims, error) {
	claims := new(JwtClaims)

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.cfg.Secret), nil
	}, jwt.WithExpirationRequired(), jwt.WithIssuedAt())

	return claims, err
}
