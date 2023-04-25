package auth

import (
	"time"

	"github.com/HouseCham/VetMate/config"
	"github.com/golang-jwt/jwt"
)

var Config *config.Config

// ShareConfigFile is a function that shares the
// configuration setted up in main.go
func ShareConfigFile(config *config.Config) {
	Config = config
}

// GenerateJWT is a function that generates a jwt token
// with the id of the user or vet that is logged in
func GenerateJWT(id int32, jwtSecret string) (string, error) {
	// first we generate the token with the HS256 signing method and
	// the claims stablished
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	// then we generate the jwt token string with the secret found in config file
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}