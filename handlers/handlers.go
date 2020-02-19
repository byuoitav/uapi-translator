package handlers

import (
	"net/http"

	"github.com/byuoitav/uapi-translator/structs"

	"github.com/labstack/echo"
)

//Rooms

func GetRooms(c echo.Context) error {
	//Check auth?

	roomNum := c.QueryParam("room_number")
	bldgAbbr := c.QueryParam("building_abbreviation")

	var rooms []structs.Room
	return c.JSON(http.StatusOK, rooms)
}

func GetRoomByID(c echo.Context) error {
	roomId := c.Param("room_id")

	var room structs.Room
	return c.JSON(http.StatusOK, room)
}

func GetRoomDevices(c echo.Context) error {
	roomId := c.Param("room_id")

	var devices structs.RoomDevices
	return c.JSON(http.StatusOK, devices)
}

//Devices

func GetDevices(c echo.Context) error {

	roomNum := c.QueryParam("room_number")
	bldgAbbr := c.QueryParam("building_abbreviation")
	deviceType := c.QueryParam("av_device_type")

	var devices []structs.Device
	return c.JSON(http.StatusOK, devices)
}

func GetDeviceByID(c echo.Context) error {
	deviceId := c.Param("av_device_id")

	var device structs.Device
	return c.JSON(http.StatusOK, device)
}

func GetDeviceProperties(c echo.Context) error {
	deviceId := c.Param("av_device_id")

	var deviceProperties []structs.DeviceProperty
	return c.JSON(http.StatusOK, deviceProperties)
}

func GetDeviceState(c echo.Context) error {
	deviceId := c.Param("av_device_id")

	var deviceStateAttrs []structs.DeviceStateAttribute
	return c.JSON(http.StatusOK, deviceStateAttrs)
}

//Inputs

func GetInputs(c echo.Context) error {

	roomNum := c.QueryParam("room_number")
	bldgAbbr := c.QueryParam("building_abbreviation")

	var inputs []structs.Input
	return c.JSON(http.StatusOK, inputs)
}

func GetInputByID(c echo.Context) error {
	deviceId := c.Param("av_device_id")

	var input structs.Input
	return c.JSON(http.StatusOK, input)
}

//Displays

func GetDisplays(c echo.Context) error {

	roomNum := c.QueryParam("room_number")
	bldgAbbr := c.QueryParam("building_abbreviation")

	var displays []structs.Display
	return c.JSON(http.StatusOK, displays)
}

func GetDisplayByID(c echo.Context) error {
	displayId := c.Param("av_display_id")

	var display structs.Display
	return c.JSON(http.StatusOK, display)
}

func GetDisplayConfig(c echo.Context) error {
	displayId := c.Param("av_display_id")

	var displayConfig structs.DisplayConfig
	return c.JSON(http.StatusOK, displayConfig)
}

func GetDisplayState(c echo.Context) error {
	displayId := c.Param("av_display_id")

	var displayState structs.DisplayState
	return c.JSON(http.StatusOK, displayState)
}

//Audio Outputs

func GetAudioOutputs(c echo.Context) error {

	roomNum := c.QueryParam("room_number")
	bldgAbbr := c.QueryParam("building_abbreviation")
	deviceType := c.QueryParam("av_device_type")

	var outputs []structs.AudioOutput
	return c.JSON(http.StatusOK, outputs)
}

func GetAudioOutputByID(c echo.Context) error {
	outputId := c.Param("av_audio_output_id")

	var output structs.AudioOutput
	return c.JSON(http.StatusOK, output)
}

func GetAudioOutputState(c echo.Context) error {
	outputId := c.Param("av_audio_output_id")

	var outputState structs.AudioOutputState
	return c.JSON(http.StatusOK, outputState)
}
