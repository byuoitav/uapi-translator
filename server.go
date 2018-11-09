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

	// State Endpoints
	router.GET("/av_rooms/av_state/:roomID", handlers.GetState)
	router.PUT("/av_rooms/av_state/:roomID", handlers.SetState)

	// Config Endpoints
	router.GET("/av_rooms/av_config/:roomID", handlers.GetConfig)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
