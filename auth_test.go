package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/softika/auth"
)

func TestAuthHandler(t *testing.T) {
	t.Parallel()

	cfg := auth.Config{
		Secret:   "test-secret",
		TokenExp: time.Hour,
	}

	tests := []struct {
		name     string
		token    func(*testing.T) string
		opts     []auth.Option
		wantCode int
	}{
		{
			name: "valid token",
			token: func(t *testing.T) string {
				valid, err := generateToken(cfg.TokenExp, cfg.Secret, jwt.SigningMethodHS256)
				if err != nil {
					t.Fatal(err)
				}
				return valid
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "invalid token",
			token:    func(t *testing.T) string { return "invalid-token" },
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "empty token",
			token:    func(t *testing.T) string { return "" },
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "expired token",
			token: func(t *testing.T) string {
				expired, err := generateToken(-time.Hour, cfg.Secret, jwt.SigningMethodHS256)
				if err != nil {
					t.Fatal(err)
				}
				return expired
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "invalid signature",
			token: func(t *testing.T) string {
				invalid, err := generateToken(cfg.TokenExp, "fake-secret", jwt.SigningMethodHS256)
				if err != nil {
					t.Fatal(err)
				}
				return invalid
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "invalid admin token",
			token: func(t *testing.T) string {
				valid, err := generateToken(cfg.TokenExp, cfg.Secret, jwt.SigningMethodHS256)
				if err != nil {
					t.Fatal(err)
				}
				return valid
			},
			opts:     []auth.Option{auth.OnlyAdmin()},
			wantCode: http.StatusForbidden,
		},
		{
			name: "valid admin token",
			token: func(t *testing.T) string {
				valid, err := generateToken(cfg.TokenExp, cfg.Secret, jwt.SigningMethodHS256, "ADMIN")
				if err != nil {
					t.Fatal(err)
				}
				return valid
			},
			opts:     []auth.Option{auth.OnlyAdmin()},
			wantCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// http request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token(t))

			// http response recorder
			w := httptest.NewRecorder()

			a := auth.New(cfg, tt.opts...)

			// wrap the test handler with the middleware
			handler := a.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			// serve http
			handler.ServeHTTP(w, req)

			// assert
			if status := w.Code; status != tt.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantCode)
			}
		})
	}
}

func generateToken(exp time.Duration, secret string, sig jwt.SigningMethod, roles ...string) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"email": "test@email.com",
		"sub":   "8bde8e50-09b3-4f70-9988-87377f791c91",
		"exp":   now.Add(exp).Unix(),
		"iat":   now.Unix(),
		"roles": roles,
	}

	token := jwt.NewWithClaims(sig, claims)

	return token.SignedString([]byte(secret))
}
