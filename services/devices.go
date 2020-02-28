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

	if roomNum != "" && bldgAbbr != "" {
		roomID := fmt.Sprintf("%s-%s", bldgAbbr, roomNum)
		dbDevices, err := requestDeviceByRoom(roomID, devType)
		if err != nil {
			return nil, err
		}
	} else if roomNum != "" {
		dbDevices, err := requestDeviceByRoomNum(roomNum, devType)
		if err != nil {
			return nil, err
		}
	} else if bldgAbbr != "" {
		dbDevices, err := requestDeviceByBuilding(bldgAbbr, devType)
		if err != nil {
			return nil, err
		}
	} else {
		dbDevices, err := requestAllDevices(devType)
		if err != nil {
			return nil, err
		}
	}

	var devices []models.Device
	for _, dev := range dbDevices {
		next := &models.Device{}
		devices = append(devices, next)
	}
	return devices, nil
}

func GetDeviceByID(deviceID string) (models.Device, error) {
	devices, err := requestDeviceByID(deviceID)
	if err != nil {
		return nil, err
	}

	s := strings.Split(deviceID, "-")
	device := &models.Device{
		deviceID:   devices[0].ID,
		deviceName: devices[0].Name,
		deviceType: devices[0].Type.ID,
		bldgAbbr:   s[0],
		roomNum:    s[1],
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

}

func requestDeviceByRoomNum(roomNum, deviceType string) ([]models.DeviceDB, error) {

}

func requestDeviceByBuilding(bldgAbbr, deviceType string) ([]models.DeviceDB, error) {

}

func requestAllDevices(deviceType string) ([]models.DeviceDB, error) {
	var query models.PrefixQuery
	query.Limit = 30 //Todo: get a definite answer on the limit
	if deviceType == nil {
		query.Selector.ID.GT = "\x00"
	} else {

	}

	url := fmt.Sprintf("%s/devices/_find", os.Getenv("DB_ADDRESS"))

	devices, err := requestDeviceSearch(url, "POST", query)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func requestDeviceSearch(url, method string, query models.PrefixQuery) ([]models.DeviceDB, error) {
	var body []byte
	if query != nil {
		body, err := json.Marshal(query)
		if err != nil {
			return nil, err
		}
	}

	var resp models.DeviceResponse
	err := couch.MakeRequest(method, url, "application/json", body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Docs, nil
}
