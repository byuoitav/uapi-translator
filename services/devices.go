package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/byuoitav/uapi-translator/couch"
	"github.com/byuoitav/uapi-translator/models"
)

func GetDevices(roomNum, bldgAbbr, devType string) ([]models.Device, error) {

	var dbDevices []models.DeviceDB
	var err error

	if roomNum != "" && bldgAbbr != "" {
		roomID := fmt.Sprintf("%s-%s", bldgAbbr, roomNum)
		dbDevices, err = requestDeviceByRoom(roomID, devType)
		if err != nil {
			return nil, err
		}
	} else if roomNum != "" {
		dbDevices, err = requestDeviceByRoomNum(roomNum, devType)
		if err != nil {
			return nil, err
		}
	} else if bldgAbbr != "" {
		dbDevices, err = requestDeviceByBuilding(bldgAbbr, devType)
		if err != nil {
			return nil, err
		}
	} else {
		dbDevices, err = requestAllDevices(devType)
		if err != nil {
			return nil, err
		}
	}

	var devices []models.Device
	if dbDevices == nil {
		return nil, fmt.Errorf("No devices")
	}
	for _, dev := range dbDevices {
		s := strings.Split(dev.ID, "-")
		next := models.Device{
			DeviceID:   dev.ID,
			DeviceName: dev.Name,
			DeviceType: dev.Type.ID,
			BldgAbbr:   s[0],
			RoomNum:    s[1],
		}
		devices = append(devices, next)
	}
	return devices, nil
}

func GetDeviceByID(deviceID string) (*models.Device, error) {
	devices, err := requestDeviceByID(deviceID)
	if err != nil {
		return nil, err
	}

	s := strings.Split(devices[0].ID, "-")
	device := &models.Device{
		DeviceID:   devices[0].ID,
		DeviceName: devices[0].Name,
		DeviceType: devices[0].Type.ID,
		BldgAbbr:   s[0],
		RoomNum:    s[1],
	}
	return device, nil
}

func requestDeviceByID(deviceID string) ([]models.DeviceDB, error) {
	url := fmt.Sprintf("%s/devices/%s", os.Getenv("DB_ADDRESS"), deviceID)

	var resp models.DeviceDB
	err := couch.MakeRequest("GET", url, "", nil, &resp)
	if err != nil {
		return nil, err
	}

	var devices []models.DeviceDB
	return append(devices, resp), nil
}

func requestDeviceByRoom(roomID, deviceType string) ([]models.DeviceDB, error) {
	url := fmt.Sprintf("%s/devices/_find", os.Getenv("DB_ADDRESS"))
	var query models.CouchQuery

	if deviceType != "" {
		query.Selector.DevType = &models.DeviceTypeQuery{
			ID: &models.CouchSearch{
				Regex: deviceType,
			},
		}
	}
	query.Limit = 1000
	query.Selector.ID.Regex = fmt.Sprintf("%s-", roomID)

	devices, err := requestDeviceSearch(url, "POST", &query)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func requestDeviceByRoomNum(roomNum, deviceType string) ([]models.DeviceDB, error) {
	url := fmt.Sprintf("%s/devices/_find", os.Getenv("DB_ADDRESS"))
	var query models.CouchQuery

	if deviceType != "" {
		query.Selector.DevType = &models.DeviceTypeQuery{
			ID: &models.CouchSearch{
				Regex: deviceType,
			},
		}
	}
	query.Limit = 1000
	query.Selector.ID.Regex = fmt.Sprintf("%s-", roomNum)

	devices, err := requestDeviceSearch(url, "POST", &query)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func requestDeviceByBuilding(bldgAbbr, deviceType string) ([]models.DeviceDB, error) {
	url := fmt.Sprintf("%s/devices/_find", os.Getenv("DB_ADDRESS"))
	var query models.CouchQuery

	if deviceType != "" {
		query.Selector.DevType = &models.DeviceTypeQuery{
			ID: &models.CouchSearch{
				Regex: deviceType,
			},
		}
	}
	query.Limit = 30 //Todo: get a definite answer on the limit
	query.Selector.ID.Regex = bldgAbbr

	devices, err := requestDeviceSearch(url, "POST", &query)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func requestAllDevices(deviceType string) ([]models.DeviceDB, error) {
	url := fmt.Sprintf("%s/devices/_find", os.Getenv("DB_ADDRESS"))
	var query models.CouchQuery

	if deviceType != "" {
		query.Selector.DevType = &models.DeviceTypeQuery{
			ID: &models.CouchSearch{
				Regex: deviceType,
			},
		}
	}
	query.Limit = 30 //Todo: get a definite answer on the limit
	query.Selector.ID.GT = "\x00"

	devices, err := requestDeviceSearch(url, "POST", &query)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func requestDeviceSearch(url, method string, query interface{}) ([]models.DeviceDB, error) {
	var body []byte
	var err error
	if query != nil {
		body, err = json.Marshal(query)
		if err != nil {
			return nil, err
		}
	}

	var resp models.DeviceResponse
	err = couch.MakeRequest(method, url, "application/json", body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Docs, nil
}
