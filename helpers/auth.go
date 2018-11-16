package helpers

import (
	"os"

	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/endpoint-authorization-controller/base"
	"github.com/labstack/echo"
)

// fieldSetRoles is a mapping of field sets to their required roles
var fieldSetRoles = map[string]string{
	Basic:    "none",
	SetState: "write-state",
	GetState: "read-state",
	Config:   "read-config",
}

// AuthenticatedByJWT scans the JWT token and processes to see if the user is granted access to the resource.
func AuthenticatedByJWT(context echo.Context, fieldSet string) (bool, error) {
	roomID := context.Param("roomID")

	role := getFieldSetRole(fieldSet, context.Request().Method)

	log.L.Debug("Arg matey, what's yer JWT like?")
	token := context.Request().Header.Get("X-jwt-assertion")
	log.L.Debug(token)

	if len(token) > 0 { // Proceed if we found a token
		// build a base.Request object with the information
		jwtRequest := base.Request{
			AccessKey: "",
			UserInformation: base.UserInformation{
				ResourceID: roomID,
				CommonInfo: base.CommonInfo{
					ID:           "",
					AuthMethod:   "wso2",
					ResourceType: "room",
					Data:         []byte(token),
				},
			},
		}

		log.L.Debug("Time to see if yer full o' barnacles or not!")
		// create an HTTP request
		req, err := jsonhttp.CreateRequest("POST", os.Getenv("ENDPOINT_AUTHORIZATION_URL"), jwtRequest, nil)
		if err != nil {
			log.L.Errorf("failed to make request - %s", err.Error())
			return false, err
		}

		// send it to the EAC
		var resp base.Response
		err = jsonhttp.ExecuteRequest(req, &resp, 60)
		if err != nil {
			log.L.Errorf("something went wrong - %s", err.Error())
		}

		if role == fieldSetRoles[Basic] {
			return true, nil
		}

		// read the response to see if the token is valid
		if resp.Permissions == nil {
			log.L.Errorf("no permissions found, to Davy Jones' locker with ye!")
			return false, nil
		}

		log.L.Debug("Searching the map for treasure....")
		for _, r := range resp.Permissions[roomID] {
			if r == role {
				return true, nil
			}
		}
	}

	log.L.Debug("Stowaway! The Cap'n'll be hear'n about this!")
	return false, nil
}

func getFieldSetRole(fs string, method string) string {
	switch fs {
	case State:
		if method == "GET" {
			return fieldSetRoles[GetState]
		}

		return fieldSetRoles[SetState]

	default:
		return fieldSetRoles[fs]
	}
}
