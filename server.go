package main

import (
	"net/http"

	"github.com/byuoitav/common"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/uapi-translator/handlers"
)

func main() {
	log.L.Info("Brwaap! Skittles the parrot!")

	port := ":9101"

	router := common.NewRouter()

	router.GET("/av_rooms/av_state/:roomID", handlers.GetState)
	router.PUT("/av_rooms/av_state/:roomID", handlers.SetState)
	router.GET("/av_rooms/av_config/:roomID", handlers.GetConfig)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
