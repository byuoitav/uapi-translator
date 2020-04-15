package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/byuoitav/uapi-translator/handlers"
	"github.com/byuoitav/uapi-translator/log"
	"github.com/byuoitav/uapi-translator/middleware"
	"github.com/labstack/echo"
	"github.com/spf13/pflag"
)

func main() {
	var port int
	var logLevel int
	var opaURL string
	var opaToken string
	var disableAuth bool

	pflag.IntVarP(&port, "port", "p", 9101, "port to run the server on")
	pflag.IntVarP(&logLevel, "log-level", "l", 2, "level of logging wanted. 1=DEBUG, 2=INFO, 3=WARN, 4=ERROR, 5=PANIC")
	pflag.StringVar(&opaURL, "opa-url", "", "URL where the OPA server can be found")
	pflag.StringVar(&opaToken, "opa-token", "", "token to use when calling OPA")
	pflag.BoolVar(&disableAuth, "disable-auth", false, "disables authz/n checks")
	pflag.Parse()

	setLog := func(level int) error {
		switch level {
		case 1:
			fmt.Printf("\nSetting log level to *debug*\n\n")
			log.Config.Level.SetLevel(zap.DebugLevel)
		case 2:
			fmt.Printf("\nSetting log level to *info*\n\n")
			log.Config.Level.SetLevel(zap.InfoLevel)
		case 3:
			fmt.Printf("\nSetting log level to *warn*\n\n")
			log.Config.Level.SetLevel(zap.WarnLevel)
		case 4:
			fmt.Printf("\nSetting log level to *error*\n\n")
			log.Config.Level.SetLevel(zap.ErrorLevel)
		case 5:
			fmt.Printf("\nSetting log level to *panic*\n\n")
			log.Config.Level.SetLevel(zap.PanicLevel)
		default:
			return errors.New("invalid log level: must be [1-4]")
		}

		return nil
	}

	// set the initial log level
	if err := setLog(logLevel); err != nil {
		log.Log.Fatal("unable to set log level", zap.Error(err), zap.Int("got", logLevel))
	}

	router := echo.New()

	// If authz/n hasn't been disabled
	if !disableAuth {
		if opaURL == "" {
			log.Log.Errorf("No OPA URL was set, but authz has not been disabled")
			os.Exit(1)
		}
		opaClient := middleware.OPAClient{
			URL:   opaURL,
			Token: opaToken,
		}

		router.Use(opaClient.Authorize)
	}

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

	// Set log level
	router.GET("/log/:level", func(c echo.Context) error {
		level, err := strconv.Atoi(c.Param("level"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if err := setLog(level); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.String(http.StatusOK, fmt.Sprintf("Set log level to %v", level))
	})

	addr := fmt.Sprintf(":%d", port)
	err := router.Start(addr)
	if err != nil {
		log.Log.Fatal("failed to start server", zap.Error(err))
	}
}
