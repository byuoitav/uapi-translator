package services

import (
	"fmt"

	"github.com/byuoitav/uapi-translator/models"
)

func GetRooms(roomNum, bldAbbr string) ([]models.Room, error) {

	var dbRooms interface{} //Todo: make a struct to receive couch rooms

	if roomNum != "" && bldAbbr != "" {
		//Both
		dbRooms, err := requestRoomByID(fmt.Sprintf("%s-%s", bldAbbr, roomNum))
		if err != nil {
			//Error getting room from database
			return nil, err
		}

	} else if roomNum != "" {
		//Just roomNum
		dbRooms, err := requestRoomByNumber(roomNum)
		if err != nil {
			return nil, err
		}
	} else if bldAbbr != "" {
		//Just bldAbbr - get rooms by building
		dbRooms, err := requestRoomByBuilding(bldAbbr)
		if err != nil {
			//Error getting rooms by building from database
			return nil, err
		}
	} else {
		//None - get all rooms
		dbRooms, err := requestAllRooms()
		if err != nil {
			//Error getting all rooms from database
			return nil, err
		}
	}

	//Translate to []models.Room
	var rooms []models.Room
	// for _, rm := range dbRooms {
	// 	s := strings.Split(rm.ID, "-")
	// 	next := &models.Room{
	// 		roomID: rm.ID,
	// 		roomNum: s[1],
	// 		bldAbbr: s[0]
	// 	}
	// 	rooms = append(rooms, next)
	// }
	return rooms, nil
}

//Request all rooms with a limit of 30?
func requestAllRooms() {
	// couch.RequestRoom
	// Post, url/room/_find, application/json, IDPrefixQuery??? - query limit: 30?
}

func requestRoomByBuilding() {}
func requestRoomByNumber()   {}
func requestRoomByID()       {}
