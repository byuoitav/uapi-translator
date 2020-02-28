package handlers

import (
	"net/http"

	"github.com/byuoitav/uapi-translator/models"
	"github.com/byuoitav/uapi-translator/services"

	"github.com/labstack/echo"
)

//Rooms

func GetRooms(c echo.Context) error {
	//Check auth?

	roomNum := c.QueryParam("room_number")
	bldgAbbr := c.QueryParam("building_abbreviation")

	rooms, err := services.GetRooms(roomNum, bldgAbbr)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, rooms)
}

func GetRoomByID(c echo.Context) error {
	roomId := c.Param("room_id")

	room, err := services.GetRoomByID(roomId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, room)
}

func GetRoomDevices(c echo.Context) error {
	// roomId := c.Param("room_id")

	var devices models.RoomDevices
	return c.JSON(http.StatusOK, devices)
}

//Devices

func GetDevices(c echo.Context) error {

	// roomNum := c.QueryParam("room_number")
	// bldgAbbr := c.QueryParam("building_abbreviation")
	// deviceType := c.QueryParam("av_device_type")

	var devices []models.Device
	return c.JSON(http.StatusOK, devices)
}

func GetDeviceByID(c echo.Context) error {
	// deviceId := c.Param("av_device_id")

	var device models.Device
	return c.JSON(http.StatusOK, device)
}

func GetDeviceProperties(c echo.Context) error {
	// deviceId := c.Param("av_device_id")

	var deviceProperties []models.DeviceProperty
	return c.JSON(http.StatusOK, deviceProperties)
}

func GetDeviceState(c echo.Context) error {
	// deviceId := c.Param("av_device_id")

	var deviceStateAttrs []models.DeviceStateAttribute
	return c.JSON(http.StatusOK, deviceStateAttrs)
}

//Inputs

func GetInputs(c echo.Context) error {

	// roomNum := c.QueryParam("room_number")
	// bldgAbbr := c.QueryParam("building_abbreviation")

	var inputs []models.Input
	return c.JSON(http.StatusOK, inputs)
}

func GetInputByID(c echo.Context) error {
	// deviceId := c.Param("av_device_id")

	var input models.Input
	return c.JSON(http.StatusOK, input)
}

//Displays

func GetDisplays(c echo.Context) error {

	// roomNum := c.QueryParam("room_number")
	// bldgAbbr := c.QueryParam("building_abbreviation")

	var displays []models.Display
	return c.JSON(http.StatusOK, displays)
}

func GetDisplayByID(c echo.Context) error {
	// displayId := c.Param("av_display_id")

	var display models.Display
	return c.JSON(http.StatusOK, display)
}

func GetDisplayConfig(c echo.Context) error {
	// displayId := c.Param("av_display_id")

	var displayConfig models.DisplayConfig
	return c.JSON(http.StatusOK, displayConfig)
}

func GetDisplayState(c echo.Context) error {
	// displayId := c.Param("av_display_id")

	var displayState models.DisplayState
	return c.JSON(http.StatusOK, displayState)
}

//Audio Outputs

func GetAudioOutputs(c echo.Context) error {

	// roomNum := c.QueryParam("room_number")
	// bldgAbbr := c.QueryParam("building_abbreviation")
	// deviceType := c.QueryParam("av_device_type")

	var outputs []models.AudioOutput
	return c.JSON(http.StatusOK, outputs)
}

func GetAudioOutputByID(c echo.Context) error {
	// outputId := c.Param("av_audio_output_id")

	var output models.AudioOutput
	return c.JSON(http.StatusOK, output)
}

func GetAudioOutputState(c echo.Context) error {
	// outputId := c.Param("av_audio_output_id")

	var outputState models.AudioOutputState
	return c.JSON(http.StatusOK, outputState)
}
