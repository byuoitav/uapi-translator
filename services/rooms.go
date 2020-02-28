package services

import (
	"fmt"
	"os"

	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/uapi-translator/couch"
	"github.com/byuoitav/uapi-translator/models"
)

func GetRooms(roomNum, bldAbbr string) ([]models.Room, error) {

	var dbRooms []models.RoomDB

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

	//Translate to []models.Room from []models.RoomDB
	var rooms []models.Room
	for _, rm := range dbRooms {
		s := strings.Split(rm.ID, "-")
		next := &models.Room{
			roomID: rm.ID,
			roomNum: s[1],
			bldAbbr: s[0]
		}
		rooms = append(rooms, next)
	}
	return rooms, nil
}

func GetRoomByID(roomID string) (models.Room, error) {
	rooms, err := requestRoomByID(roomID)
	if err != nil {
		return nil, err
	}
	return rooms[0], nil
}

func GetRoomDevices(roomID string) ([]models.RoomDevices, error) {
	rooms, err := requestRoomByID(roomID)
	if err != nil {
		return nil, err
	}

	//Get devices????
	
	var devices []models.Device
	for _, d := range rooms[0].Devices {
		s := strings.Split(rooms[0].ID, "-")
		next := &models.Device {
			deviceID: d.ID,
			deviceName: d.Name,
			deviceType: d.Type.ID,
			bldgAbbr: s[0],
			roomNum: s[1],
		}
	}
	return devices, nil
}

func requestRoomByID(roomID string) ([]models.RoomDB, error) {
	url := fmt.Sprintf("%s/rooms/%s", os.Getenv("DB_ADDRESS"), roomID)

	var resp models.RoomDB
	err = couch.MakeRequest(method, url, "application/json", body, &resp)
	if err != nil {
		return nil, err
	}
	
	var rooms []models.RoomDB
	return append(rooms, resp), nil
}

func requestRoomByNumber(roomNum string) ([]models.RoomDB, error) {
	var query models.PrefixQuery
	//Todo: search for rooms with roomNum, regex?
	query.Limit = 1000

	url := fmt.Sprintf("%s/rooms/_find", os.Getenv("DB_ADDRESS"))

	rooms, err := requestRoomSearch(url, "POST", query)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func requestRoomByBuilding(bldAbbr string) ([]models.RoomDB, error) {
	var query models.PrefixQuery
	query.Selector.ID.Regex = bldAbbr
	query.Limit = 1000

	url := fmt.Sprintf("%s/rooms/_find", os.Getenv("DB_ADDRESS"))

	rooms, err := requestRoomSearch(url, "POST", query)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

//Request all rooms with a limit of 30?
func requestAllRooms() ([]models.RoomDB, error) {
	// Post, url/room/_find, application/json, IDPrefixQuery??? - query limit: 30?
	var query models.PrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 30

	url := fmt.Sprintf("%s/rooms/_find", os.Getenv("DB_ADDRESS"))

	rooms, err := requestRoomSearch(url, "POST", query)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func requestRoomSearch(url, method string, query models.PrefixQuery) ([]models.RoomDB, error) {
	var body []byte
	if query != nil {
		body, err := json.Marshal(query)
		if err != nil {
			return nil, err
		}
	}
	// call makeRequest
	var resp []models.RoomResponse
	err = couch.MakeRequest(method, url, "application/json", body, &resp)
	if err != nil {
		return nil, err
	}
	// translate response body
	return resp.Docs, nil
}
