package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/byuoitav/uapi-translator/couch"
	"github.com/byuoitav/uapi-translator/models"
)

func GetDevices(roomNum, bldgAbbr, devType string) ([]models.Device, error) {
	url := fmt.Sprintf("%s/devices/_find", os.Getenv("DB_ADDRESS"))
	var query models.DeviceQuery

	if devType != "" {
		query.Selector.DevType = &models.DeviceTypeQuery{
			ID: &models.CouchSearch{
				Regex: devType,
			},
		}
	}

	if roomNum != "" && bldgAbbr != "" {
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-%s-", bldgAbbr, roomNum)
	} else if roomNum != "" {
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-", roomNum)
	} else if bldgAbbr != "" {
		query.Limit = 30 //Todo: get a definite answer on the limit
		query.Selector.ID.Regex = bldgAbbr
	} else {
		query.Limit = 30 //Todo: get a definite answer on the limit
		query.Selector.ID.GT = "\x00"
	}

	var resp models.DeviceResponse
	err := couch.DBSearch(url, "POST", &query, &resp)
	if err != nil {
		return nil, err
	}

	var devices []models.Device
	if resp.Docs == nil {
		return nil, fmt.Errorf("No devices")
	}
	for _, dev := range resp.Docs {
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
