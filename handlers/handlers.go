package handlers

import "github.com/labstack/echo"

//Rooms
func GetRooms(c echo.Context) error       {}
func GetRoomByID(c echo.Context) error    {}
func GetRoomDevices(c echo.Context) error {}

//Devices
func GetDevices(c echo.Context) error          {}
func GetDeviceByID(c echo.Context) error       {}
func GetDeviceProperties(c echo.Context) error {}
func GetDeviceState(c echo.Context) error      {}

//Inputs
func GetInputs(c echo.Context) error    {}
func GetInputByID(c echo.Context) error {}

//Displays
func GetDisplays(c echo.Context) error      {}
func GetDisplayByID(c echo.Context) error   {}
func GetDisplayConfig(c echo.Context) error {}
func GetDisplayState(c echo.Context) error  {}

//Audio Outputs
func GetAudioOutputs(c echo.Context) error     {}
func GetAudioOutputByID(c echo.Context) error  {}
func GetAudioOutputState(c echo.Context) error {}
