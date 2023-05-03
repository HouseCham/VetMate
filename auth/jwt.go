package auth

import (
	"strconv"
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
// isVet is a boolean that indicates if the user is a vet
// or not
func GenerateJWT(id int32, isVet bool) (string, error) {
	// first we generate the token with the HS256 signing method and
	// the claims stablished
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(id)),
		"isVet": isVet,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	// then we generate the jwt token string with the secret found in config file
	tokenString, err := token.SignedString([]byte(Config.DevConfiguration.Jwt.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}