package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/v2/auth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// AuthenticatedByJWT scans the JWT token and processes to see if the user is granted access to the resource.
func AuthenticatedByJWT(context echo.Context, role string) (bool, error) {
	roomID := context.Param("room")
	token := context.Request().Header.Get("X-jwt-assertion")

	if len(token) > 0 { // Proceed if we found a token
		valid, err := validate(token) // Validate the existing token
		if err != nil {
			log.L.Debug("Invalid WSO2 information")
			return false, err
		}

		if valid {
			log.L.Debug("WSO2 validated successfully")

			// check user's roles in our database
			// TODO: get the real user and access key and stuff...
			ok, err := auth.CheckRolesForUser("user", "accessKey", role, roomID, "room")
			if err != nil {
				return false, err
			}

			if !ok {
				return false, nil
			}

			return true, nil
		}
	}

	return false, nil
}

type keys struct {
	Keys []struct {
		E   string   `json:"e"`
		Kty string   `json:"kty"`
		Use string   `json:"use"`
		Kid string   `json:"kid"`
		N   string   `json:"n"`
		X5C []string `json:"x5c"`
	} `json:"keys"`
}

func validate(token string) (bool, error) {
	parsedToken, err := jwt.Parse(token, func(parsedToken *jwt.Token) (interface{}, error) {
		if parsedToken.Method.Alg() != "RS256" { // Check that our keys are signed with RS256 as expected (https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/)
			return nil, fmt.Errorf("Unexpected signing method: %v", parsedToken.Header["alg"]) // This error never gets returned to the user but may be useful for debugging/logging at some point
		}

		// Look up key
		key, err := lookupSigningKey()
		if err != nil {
			return nil, err
		}

		// Unpack key from PEM encoded PKCS8
		return jwt.ParseRSAPublicKeyFromPEM(key)
	})

	log.L.Debugf("%v", parsedToken)

	if parsedToken.Valid {
		return true, nil
	} else if validationError, ok := err.(*jwt.ValidationError); ok {
		if validationError.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, errors.New("Authorization token is malformed")
		} else if validationError.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return false, errors.New("Authorization token is expired")
		}
	}

	return false, errors.New("WSO2 JWT token not authorized")
}

func lookupSigningKey() ([]byte, error) {
	response, err := http.Get("https://api.byu.edu/.well-known/byucerts")
	if err != nil {
		return nil, err
	}

	allKeys := keys{}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &allKeys)
	if err != nil {
		return nil, err
	}

	certificate := "-----BEGIN CERTIFICATE-----\n" + allKeys.Keys[0].X5C[0] + "\n-----END CERTIFICATE-----"
	log.L.Debug(certificate)
	return []byte(certificate), nil
}
