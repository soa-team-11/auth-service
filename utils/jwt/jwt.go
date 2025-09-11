package jwt

import (
	"github.com/go-chi/jwtauth/v5"

	"github.com/soa-team-11/auth-service/utils"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(utils.Getenv("JWT_SECRET", "soa_team_11")), nil)
}

func GetTokenAuth() *jwtauth.JWTAuth {
	return tokenAuth
}
