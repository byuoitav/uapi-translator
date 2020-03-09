package services

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/byuoitav/scheduler/log"
	"github.com/byuoitav/uapi-translator/db"
	"github.com/byuoitav/uapi-translator/models"
)

func GetRooms(roomNum, bldgAbbr string) ([]models.Room, error) {
	url := fmt.Sprintf("%s/rooms/_find", os.Getenv("DB_ADDRESS"))
	var query models.RoomQuery

	if roomNum != "" && bldgAbbr != "" {
		log.P.Info("searching rooms by room number and building abbreviation", zap.String("roomNum", roomNum), zap.String("bldgAbbr", bldgAbbr))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-%s$", bldgAbbr, roomNum)
	} else if roomNum != "" {
		log.P.Info("searching rooms by room number", zap.String("roomNum", roomNum))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("-%s$", roomNum)
	} else if bldgAbbr != "" {
		log.P.Info("searching rooms by building abbreviation", zap.String("bldgAbbr", bldgAbbr))
		query.Limit = 30 //Todo: get a definite answer on the limit
		query.Selector.ID.Regex = bldgAbbr
	} else {
		log.P.Info("getting all rooms")
		query.Limit = 30 //Todo: get a definite answer on the limit
		query.Selector.ID.GT = "\x00"
	}

	var resp models.RoomResponse
	err := db.DBSearch(url, "POST", &query, &resp)
	if err != nil {
		log.P.Error("failed to search for rooms in database", zap.Error(err))
		return nil, err
	}

	var rooms []models.Room
	if resp.Docs == nil {
		log.P.Info("no rooms resulted from query")
		return nil, fmt.Errorf("No rooms exist under the provided search criteria")
	}
	for _, rm := range resp.Docs {
		s := strings.Split(rm.ID, "-")
		next := models.Room{
			RoomID:   rm.ID,
			RoomNum:  s[1],
			BldgAbbr: s[0],
		}
		rooms = append(rooms, next)
	}
	return rooms, nil
}

// func GetRoomDevices(roomID string) ([]models.RoomDevices, error) {
// 	rooms, err := requestRoomByID(roomID)
// 	if err != nil {
// 		return nil, err
// 	} else if rooms == nil {
// 		return nil, nil //Return error stating no rooms found?
// 	}

//	// Get devices????

// 	var devices []models.RoomDevices
// for _, d := range rooms[0].Devices {
// 	s := strings.Split(rooms[0].ID, "-")
// 	next := &models.Device{
// 		deviceID:   d.ID,
// 		deviceName: d.Name,
// 		deviceType: d.Type.ID,
// 		bldgAbbr:   s[0],
// 		roomNum:    s[1],
// 	}
// }
// 	return devices, nil
// }
