package auth

import (
	"slices"

	"github.com/golang-jwt/jwt/v5"
)

type key string

const CtxKey key = "AUTH_CTX"

type Context struct {
	IsAdmin bool
	UserId  string
	Email   string
	Roles   []string
}

type JwtClaims struct {
	Admin bool
	Email string
	Roles []string
	jwt.RegisteredClaims
}

func (c JwtClaims) toContext() Context {
	isAdmin := c.Admin || slices.Contains(c.Roles, "ADMIN")
	return Context{
		IsAdmin: isAdmin,
		UserId:  c.Subject,
		Email:   c.Email,
		Roles:   c.Roles,
	}
}
