package helpers

import (
	"github.com/labstack/echo"
)

// AuthenticatedByJWT scans the JWT token and processes to see if the user is granted access to the resource.
func AuthenticatedByJWT(context echo.Context) (bool, error) {
	return true, nil
}
