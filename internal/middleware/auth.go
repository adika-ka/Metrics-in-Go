package middleware

import (
	"net/http"

	"github.com/go-chi/jwtauth"
)

var TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

func JWTMiddleware(next http.Handler) http.Handler {
	return jwtauth.Verifier(TokenAuth)(jwtauth.Authenticator(next))
}
