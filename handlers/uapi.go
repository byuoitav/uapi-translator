package handlers

import (
	"net/http"

	"github.com/byuoitav/av-api/base"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/uapi-translator/helpers"
	"github.com/labstack/echo"
)

// SetState translates the body from the UAPI into the AV-API format and forwards the request.
func SetState(context echo.Context) error {
	roomID := context.Param("roomID")

	var reqBody base.PublicRoom

	// bind the body of the request from the UAPI
	err := context.Bind(reqBody)
	if err != nil {
		log.L.Errorf("failed to bind the request body from the University API : %s", err.Error())
		return context.JSON(http.StatusBadRequest, err)
	}

	// execute the request with the new body against the AV-API
	resp, ne := AVSetState(roomID, reqBody)
	if ne != nil {
		log.L.Errorf("failed to execute request against the AV-API : %s", ne.String())
		return context.JSON(http.StatusInternalServerError, ne.String())
	}

	// translate the response back into the UAPI format
	toReturn, ne := helpers.AVtoUAPI(resp, helpers.Basic, helpers.State)
	if ne != nil {
		log.L.Errorf("failed to translate to UAPI format : %s", ne.String())
		return context.JSON(http.StatusInternalServerError, ne.String())
	}

	return context.JSON(http.StatusOK, toReturn)
}

// GetState translates the body from the UAPI into the AV-API format and forwards the request.
func GetState(context echo.Context) error {
	roomID := context.Param("roomID")

	// execute the request against the AV-API
	resp, ne := AVGetState(roomID)
	if ne != nil {
		log.L.Errorf("failed to execute request again the AV-API : %s", ne.String())
		return context.JSON(http.StatusInternalServerError, ne.String())
	}

	// translate the response into the UAPI format
	toReturn, ne := helpers.AVtoUAPI(resp, helpers.Basic, helpers.State)
	if ne != nil {
		log.L.Errorf("failed to translate to UAPI format : %s", ne.String())
		return context.JSON(http.StatusInternalServerError, ne.String())
	}

	return context.JSON(http.StatusOK, toReturn)
}

// GetConfig translates the body from the UAPI into the AV-API format and forwards the request.
func GetConfig(context echo.Context) error {
	return nil
}
