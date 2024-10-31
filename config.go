package auth

import "time"

type Config struct {
	Secret   string        `mapstructure:"secret" validate:"required"` // Secret key for JWT
	TokenExp time.Duration `mapstructure:"token_exp"`                  // Token expiration time
}
