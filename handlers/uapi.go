package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/byuoitav/common/nerr"

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

const queryKey = "?field_sets"

// GetBasic return the basic info about the resource
func GetBasic(context echo.Context) error {
	roomID := context.Param("roomID")

	values := context.QueryParams()

	var filters []string

	if values[queryKey] != nil {
		filters = strings.Split(values[queryKey][0], ",")
	}

	var fieldSets = []string{"basic"}

	fieldSets = append(fieldSets, filters...)

	toReturn, code, err := handleRequest(context, roomID, fieldSets...)
	if err != nil {
		return context.JSON(code, err)
	}

	log.L.Debugf("Woohoo! Successfully got the basic field set of %s!", roomID)
	return context.JSON(http.StatusOK, toReturn)
}

// State translates the body from the UAPI into the AV-API format and forwards the request.
func State(context echo.Context) error {
	roomID := context.Param("roomID")

	values := context.QueryParams()

	var filters []string

	if values[queryKey] != nil {
		filters = strings.Split(values[queryKey][0], ",")
	}

	var fieldSets = []string{"av_state"}

	fieldSets = append(fieldSets, filters...)

	toReturn, code, err := handleRequest(context, roomID, fieldSets...)
	if err != nil {
		return context.JSON(code, err)
	}

	log.L.Debug(color.HiMagentaString("Woohoo! Successfully set the state of %s!", roomID))
	return context.JSON(http.StatusOK, toReturn)
}

// GetConfig translates the body from the UAPI into the AV-API format and forwards the request.
func GetConfig(context echo.Context) error {
	roomID := context.Param("roomID")

	values := context.QueryParams()

	var filters []string

	if values[queryKey] != nil {
		filters = strings.Split(values[queryKey][0], ",")
	}

	var fieldSets = []string{"av_config"}

	fieldSets = append(fieldSets, filters...)

	toReturn, code, err := handleRequest(context, roomID, fieldSets...)
	if err != nil {
		return context.JSON(code, err)
	}

	log.L.Debugf("Woohoo! Successfully got the configuration of %s!", roomID)
	return context.JSON(http.StatusOK, toReturn)
}

func handleRequest(context echo.Context, roomID string, fieldSets ...string) (structs.Resource, int, error) {
	var pr base.PublicRoom
	var rr structs.ReachableRoomConfig
	var ne *nerr.E

	for _, fs := range fieldSets {
		// Check the validity of the JWT token and to see if the user has the necessary roles.
		ok, err := helpers.AuthenticatedByJWT(context, fs)
		if err != nil {
			log.L.Error(color.HiRedString("Brwaap! Error authenticating! : %s", err.Error()))
			return structs.Resource{}, http.StatusInternalServerError, fmt.Errorf("There was a problem")
		}

		if !ok {
			log.L.Error(color.HiRedString("Brwaap! The scalawag did not meet authentication checks!"))
			if fs == helpers.State {
				pr.Building = helpers.Unauthorized
			}
			if fs == helpers.Config {
				rr.ID = helpers.Unauthorized
			}
		} else {
			// Process the request with the AV-API
			switch fs {
			case helpers.Config:
				rr, ne = AVGetConfig(roomID)
				if ne != nil {
					log.L.Errorf("an error while getting the config - %s", ne.String())
					return structs.Resource{}, http.StatusInternalServerError, ne
				}
			case helpers.State:
				if context.Request().Method == "GET" {
					pr, ne = AVGetState(roomID)
					if ne != nil {
						log.L.Errorf(color.HiRedString("Brwaap! Time to walk the plank! Failed to execute request again the AV-API : %s", ne.String()))
						return structs.Resource{}, http.StatusInternalServerError, ne
					}
				} else {
					var reqBody base.PublicRoom

					// bind the body of the request from the UAPI
					err = context.Bind(&reqBody)
					if err != nil {
						log.L.Error(color.HiRedString("Brwaap! Time to walk the plank! Failed to bind the request body from the University API : %s", err.Error()))
						return structs.Resource{}, http.StatusBadRequest, err
					}

					// execute the request with the new body against the AV-API
					pr, ne = AVSetState(roomID, reqBody)
					if ne != nil {
						log.L.Errorf(color.HiRedString("Brwaap! Time to walk the plank! Failed to execute request against the AV-API : %s", ne.String()))
						return structs.Resource{}, http.StatusInternalServerError, ne
					}
				}
			}
		}
	}

	// Translate the response into UAPI format
	toReturn, ne := helpers.AVtoUAPI(roomID, pr, rr, fieldSets...)
	if ne != nil {
		log.L.Errorf(color.HiRedString("Brwaap! Time to walk the plank! Failed to translate to UAPI format : %s", ne.String()))
		return structs.Resource{}, http.StatusInternalServerError, ne
	}

	// Return the response
	return toReturn, http.StatusOK, nil
}
