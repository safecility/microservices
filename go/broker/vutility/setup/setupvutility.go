package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/broker/vutility/helpers"
	"os"
	"time"
)

func main() {
	deployment, isSet := os.LookupEnv("Deployment")
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	jwtSecretName := fmt.Sprintf("projects/%s/secrets/jwt-key/versions/1", config.ProjectName)
	sigSecret := setup.GetSecret(jwtSecretName)

	sig := hmac.New(sha256.New, []byte(sigSecret))

	hmacSecret := sig.Sum(nil)

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"service": "vutility",
		"created": now.Format(time.RFC3339),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		log.Err(err).Msg("could not generate token")
		return
	}

	fmt.Println(tokenString, sig)

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if token == nil || !token.Valid {
		log.Err(err).Msg("invalid token")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["service"], claims["created"])
	} else {
		fmt.Println(err)
	}
}
