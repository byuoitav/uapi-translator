package main

import (
	"net/http"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/uapi-translator/handlers"
	"github.com/fatih/color"
)

func main() {
	log.L.Infof("%s %s %s", color.HiCyanString("Brwaap!"), color.HiGreenString("GoGo"), color.HiYellowString("the parrot!"))

	port := ":9101"

	router := common.NewRouter()

	// Log Endpoints
	router.PUT("/log-level/:level", log.SetLogLevel)
	router.GET("/log-level", log.GetLogLevel)

	// Basic Endpoint
	router.GET("/:roomID", handlers.GetBasic)

	// State Endpoints
	router.GET("/:roomID/av_state", handlers.State)
	router.PUT("/:roomID/av_state", handlers.State)

	// Config Endpoints
	router.GET("/:roomID/av_config", handlers.GetConfig)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
