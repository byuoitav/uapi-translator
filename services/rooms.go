package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/byuoitav/uapi-translator/couch"
	"github.com/byuoitav/uapi-translator/models"
)

func GetRooms(roomNum, bldgAbbr string) ([]models.Room, error) {
	url := fmt.Sprintf("%s/rooms/_find", os.Getenv("DB_ADDRESS"))
	var query models.RoomQuery

	if roomNum != "" && bldgAbbr != "" {
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-%s$", bldgAbbr, roomNum)
	} else if roomNum != "" {
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("-%s$", roomNum)
	} else if bldgAbbr != "" {
		query.Limit = 30 //Todo: get a definite answer on the limit
		query.Selector.ID.Regex = bldgAbbr
	} else {
		query.Limit = 30 //Todo: get a definite answer on the limit
		query.Selector.ID.GT = "\x00"
	}

	var resp models.RoomResponse
	err := couch.DBSearch(url, "POST", &query, &resp)
	if err != nil {
		return nil, err
	}

	var rooms []models.Room
	if resp.Docs == nil {
		return nil, fmt.Errorf("No rooms")
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
