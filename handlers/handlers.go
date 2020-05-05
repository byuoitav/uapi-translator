package handlers

import (
	"net/http"
	"strings"

	"github.com/byuoitav/uapi-translator/log"
	"github.com/byuoitav/uapi-translator/models"
	"github.com/byuoitav/uapi-translator/services"

	"github.com/labstack/echo"
)

type Service struct {
	Services *services.Service
}

//Rooms

func (s *Service) GetRooms(c echo.Context) error {
	roomNum := c.QueryParam("room_number")
	bldgAbbr := c.QueryParam("building_abbreviation")

	rooms, err := s.Services.GetRooms(roomNum, bldgAbbr)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Infof("successfully retrieved: %d rooms", len(rooms))
	return c.JSON(http.StatusOK, rooms)
}

func (s *Service) GetRoomByID(c echo.Context) error {
	roomId := c.Param("room_id")
	parts := strings.Split(roomId, "-")

	room, err := s.Services.GetRooms(parts[1], parts[0])
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Info("successfully retrieved room by id")
	return c.JSON(http.StatusOK, room)
}

func (s *Service) GetRoomDevices(c echo.Context) error {
	roomId := c.Param("room_id")

	devices, err := s.Services.GetRoomDevices(roomId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Info("successfully retrieved room devices")
	return c.JSON(http.StatusOK, devices)
}

//Devices

func (s *Service) GetDevices(c echo.Context) error {
	roomNum := c.QueryParam("room_number")
	bldgAbbr := c.QueryParam("building_abbreviation")
	deviceType := c.QueryParam("av_device_type")

	devices, err := s.Services.GetDevices(roomNum, bldgAbbr, deviceType)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Infof("successfully retrieved: %d devices", len(devices))
	return c.JSON(http.StatusOK, devices)
}

func (s *Service) GetDeviceByID(c echo.Context) error {
	deviceId := c.Param("av_device_id")

	device, err := s.Services.GetDeviceByID(deviceId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Info("successfully retrieved device by id")
	return c.JSON(http.StatusOK, device)
}

func (s *Service) GetDeviceProperties(c echo.Context) error {
	// deviceId := c.Param("av_device_id")

	var deviceProperties []models.DeviceProperty
	return c.JSON(http.StatusOK, deviceProperties)
}

func (s *Service) GetDeviceState(c echo.Context) error {
	// deviceId := c.Param("av_device_id")

	var deviceStateAttrs []models.DeviceStateAttribute
	return c.JSON(http.StatusOK, deviceStateAttrs)
}

//Inputs

func (s *Service) GetInputs(c echo.Context) error {

	roomNum := c.QueryParam("room_number")
	bldgAbbr := c.QueryParam("building_abbreviation")

	inputs, err := s.Services.GetInputs(roomNum, bldgAbbr)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Infof("successfully retrieved: %d inputs", len(inputs))
	return c.JSON(http.StatusOK, inputs)
}

func (s *Service) GetInputByID(c echo.Context) error {
	deviceId := c.Param("av_device_id")

	input, err := s.Services.GetInputByID(deviceId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Info("successfully retrieved input by id")
	return c.JSON(http.StatusOK, input)
}

//Displays

func (s *Service) GetDisplays(c echo.Context) error {

	roomNum := c.QueryParam("room_number")
	bldgAbbr := c.QueryParam("building_abbreviation")

	displays, err := s.Services.GetDisplays(roomNum, bldgAbbr)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Infof("successfully retrieved: %d displays", len(displays))
	return c.JSON(http.StatusOK, displays)
}

func (s *Service) GetDisplayByID(c echo.Context) error {
	displayId := c.Param("av_display_id")

	display, err := s.Services.GetDisplayByID(displayId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Info("successfully retrieved display by id")
	return c.JSON(http.StatusOK, display)
}

func (s *Service) GetDisplayConfig(c echo.Context) error {
	displayId := c.Param("av_display_id")

	displayConfig, err := s.Services.GetDisplayConfig(displayId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Info("successfully retrieved display config")
	return c.JSON(http.StatusOK, displayConfig)
}

func (s *Service) GetDisplayState(c echo.Context) error {
	displayId := c.Param("av_display_id")

	displayState, err := s.Services.GetDisplayState(displayId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Info("successfully retrieved display state")
	return c.JSON(http.StatusOK, displayState)
}

//Audio Outputs

func (s *Service) GetAudioOutputs(c echo.Context) error {
	roomNum := c.QueryParam("room_number")
	bldgAbbr := c.QueryParam("building_abbreviation")
	deviceType := c.QueryParam("av_device_type")

	outputs, err := s.Services.GetAudioOutputs(roomNum, bldgAbbr, deviceType)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Infof("successfully retrieved: %d audio outputs", len(outputs))
	return c.JSON(http.StatusOK, outputs)
}

func (s *Service) GetAudioOutputByID(c echo.Context) error {
	outputId := c.Param("av_audio_output_id")

	output, err := s.Services.GetAudioOutputByID(outputId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Info("successfully retrieved audio output by id")
	return c.JSON(http.StatusOK, output)
}

func (s *Service) GetAudioOutputState(c echo.Context) error {
	outputId := c.Param("av_audio_output_id")

	outputState, err := s.Services.GetAudioOutputState(outputId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	log.Log.Info("successfully retrieved audio output state by id")
	return c.JSON(http.StatusOK, outputState)
}
