package handlers

import (
	"net/http"

	"github.com/fatih/color"

	"github.com/byuoitav/av-api/base"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/uapi-translator/helpers"
	"github.com/byuoitav/uapi-translator/structs"
	"github.com/labstack/echo"
)

var nilPubRoom = base.PublicRoom{
	Building: "",
	Room:     "",
}
var nilUAPIRoom = structs.Room{
	ID:          "",
	Name:        "",
	Description: "",
	Designation: "",
}
var nilReachableRoom = structs.ReachableRoomConfig{
	Room:              nilUAPIRoom,
	InputReachability: make(map[string][]string),
}

// SetState translates the body from the UAPI into the AV-API format and forwards the request.
func SetState(context echo.Context) error {
	roomID := context.Param("roomID")

	// check the JWT token
	ok, err := helpers.AuthenticatedByJWT(context, "write-state")
	if err != nil {
		log.L.Error(color.HiRedString("Brwaap! Error authenticating! : %s", err.Error()))
		return context.JSON(http.StatusInternalServerError, "There was a problem")
	}

	if !ok {
		log.L.Error(color.HiRedString("Brwaap! The scalawag did not meet authentication checks!"))
		return context.JSON(http.StatusForbidden, "Unauthorized")
	}

	// proceed because the JWT token checks passed
	log.L.Debug(color.HiMagentaString("Brwaap! Setting state for %s!", roomID))

	var reqBody base.PublicRoom

	// bind the body of the request from the UAPI
	err = context.Bind(&reqBody)
	if err != nil {
		log.L.Error(color.HiRedString("Brwaap! Time to walk the plank! Failed to bind the request body from the University API : %s", err.Error()))
		return context.JSON(http.StatusBadRequest, err)
	}

	// execute the request with the new body against the AV-API
	resp, ne := AVSetState(roomID, reqBody)
	if ne != nil {
		log.L.Errorf(color.HiRedString("Brwaap! Time to walk the plank! Failed to execute request against the AV-API : %s", ne.String()))
		return context.JSON(http.StatusInternalServerError, ne.String())
	}

	// translate the response back into the UAPI format
	toReturn, ne := helpers.AVtoUAPI(roomID, resp, nilReachableRoom, helpers.Basic, helpers.State)
	if ne != nil {
		log.L.Errorf(color.HiRedString("Brwaap! Time to walk the plank! Failed to translate to UAPI format : %s", ne.String()))
		return context.JSON(http.StatusInternalServerError, ne.String())
	}

	log.L.Debug(color.HiMagentaString("Woohoo! Successfully set the state of %s!", roomID))
	return context.JSON(http.StatusOK, toReturn)
}

// GetState translates the body from the UAPI into the AV-API format and forwards the request.
func GetState(context echo.Context) error {
	roomID := context.Param("roomID")

	// check the JWT token
	ok, err := helpers.AuthenticatedByJWT(context, "read-state")
	if err != nil {
		log.L.Error(color.HiRedString("Brwaap! Error authenticating! : %s", err.Error()))
		return context.JSON(http.StatusInternalServerError, "There was a problem")
	}

	if !ok {
		log.L.Error(color.HiRedString("Brwaap! The scalawag did not meet authentication checks!"))
		return context.JSON(http.StatusForbidden, "Unauthorized")
	}

	log.L.Debugf("Brwaap! Getting state for %s!", roomID)

	// execute the request against the AV-API
	resp, ne := AVGetState(roomID)
	if ne != nil {
		log.L.Errorf(color.HiRedString("Brwaap! Time to walk the plank! Failed to execute request again the AV-API : %s", ne.String()))
		return context.JSON(http.StatusInternalServerError, ne.String())
	}

	// translate the response into the UAPI format
	toReturn, ne := helpers.AVtoUAPI(roomID, resp, nilReachableRoom, helpers.Basic, helpers.State)
	if ne != nil {
		log.L.Errorf(color.HiRedString("Brwaap! Time to walk the plank! Failed to translate to UAPI format : %s", ne.String()))
		return context.JSON(http.StatusInternalServerError, ne.String())
	}

	log.L.Debugf("Woohoo! Successfully got the state of %s!", roomID)
	return context.JSON(http.StatusOK, toReturn)
}

// GetConfig translates the body from the UAPI into the AV-API format and forwards the request.
func GetConfig(context echo.Context) error {
	roomID := context.Param("roomID")

	// check the JWT token
	ok, err := helpers.AuthenticatedByJWT(context, "read-config")
	if err != nil {
		log.L.Error(color.HiRedString("Brwaap! Error authenticating! : %s", err.Error()))
		return context.JSON(http.StatusInternalServerError, "There was a problem")
	}

	if !ok {
		log.L.Error(color.HiRedString("Brwaap! The scalawag did not meet authentication checks!"))
		return context.JSON(http.StatusForbidden, "Unauthorized")
	}

	log.L.Debugf("Brwaap! Getting configuration for %s!", roomID)

	//Get the Configuration from the AV-API
	resp, ne := AVGetConfig(roomID)
	if ne != nil {
		log.L.Debugf("%s", ne.String())
		return ne
	}

	//Translate
	toReturn, ne := helpers.AVtoUAPI(roomID, nilPubRoom, resp, helpers.Config)
	if ne != nil {
		log.L.Errorf(color.HiRedString("Brwaap! Time to walk the plank! Failed to translate to UAPI format : %s", ne.String()))
		return context.JSON(http.StatusInternalServerError, ne.String())
	}

	//Return
	return context.JSON(http.StatusOK, toReturn)
}
