package handlers

import (
	"net/http"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/uapi-translator/helpers"
	"github.com/labstack/echo"
)

// SetState translates the body from the UAPI into the AV-API format and forwards the request.
func SetState(context echo.Context) error {
	roomID := context.Param("roomID")

	var reqBody structs.Resource

	// bind the body of the request from the UAPI
	err := context.Bind(reqBody)
	if err != nil {
		log.L.Errorf("failed to bind the request body from the University API : %s", err.Error())
		return context.JSON(http.StatusBadRequest, err)
	}

	// translate the Resource body to the AV-API format
	pubRoom, ne := helpers.UAPItoAV(reqBody)
	if ne != nil {
		log.L.Errorf("failed to translate to AV-API format : %s", ne.String())
		return context.JSON(http.StatusInternalServerError, ne.String())
	}

	return nil
}

// GetState translates the body from the UAPI into the AV-API format and forwards the request.
func GetState(context echo.Context) error {
	return nil
}

// GetConfig translates the body from the UAPI into the AV-API format and forwards the request.
func GetConfig(context echo.Context) error {
	return nil
}
