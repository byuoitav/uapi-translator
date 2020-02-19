package main

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/common"
	"github.com/byuoitav/uapi-translator/handlers"
	"github.com/spf13/pflag"
)

func main() {
	var port int

	pflag.IntVarP(&port, "port", "p", 9101, "port to run the server on")
	pflag.Parse()

	router := common.NewRouter()

	//Rooms
	router.GET("/rooms", handlers.GetRooms)
	router.GET("/rooms/:room_id", handlers.GetRoomByID)
	router.GET("/rooms/:room_id/devices", handlers.GetRoomDevices)

	//Devices
	router.GET("/devices", handlers.GetDevices)
	router.GET("/devices/:av_device_id", handlers.GetDeviceByID)
	router.GET("/devices/:av_device_id/properties", handlers.GetDeviceProperties)
	router.GET("/devices/:av_device_id/state", handlers.GetDeviceState)

	//Inputs
	router.GET("/inputs", handlers.GetInputs)
	router.GET("/inputs/:av_device_id", handlers.GetInputByID)

	//Displays
	router.GET("/displays", handlers.GetDisplays)
	router.GET("/displays/:av_display_id", handlers.GetDisplayByID)
	router.GET("/displays/:av_display_id/config", handlers.GetDisplayConfig)
	router.GET("/displays/:av_display_id/state", handlers.GetDisplayState)

	//Audio Outputs
	router.GET("/audio_outputs", handlers.GetAudioOutputs)
	router.GET("/audio_outputs/:av_audio_output_id", handlers.GetAudioOutputByID)
	router.GET("/audio_outputs/:av_audio_output_id/state", handlers.GetAudioOutputState)

	addr := fmt.Sprintf(":%d", port)
	err := router.StartServer(&http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1024 * 10,
	})
	if err != nil {
		// Log error
	}
}
