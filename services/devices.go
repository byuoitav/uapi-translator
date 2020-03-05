package services

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/byuoitav/scheduler/log"
	"github.com/byuoitav/uapi-translator/couch"
	"github.com/byuoitav/uapi-translator/models"
)

func GetDevices(roomNum, bldgAbbr, devType string) ([]models.Device, error) {
	url := fmt.Sprintf("%s/devices/_find", os.Getenv("DB_ADDRESS"))
	var query models.DeviceQuery

	if devType != "" {
		log.P.Info("searching with device type", zap.String("devType", devType))
		query.Selector.DevType = &models.DeviceTypeQuery{
			ID: &models.CouchSearch{
				Regex: devType,
			},
		}
	}

	if roomNum != "" && bldgAbbr != "" {
		log.P.Info("searching devices by room number and building abbreviation", zap.String("roomNum", roomNum), zap.String("bldgAbbr", bldgAbbr))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-%s-", bldgAbbr, roomNum)
	} else if roomNum != "" {
		log.P.Info("searching devices by room number", zap.String("roomNum", roomNum))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-", roomNum)
	} else if bldgAbbr != "" {
		log.P.Info("searching devices by building abbreviation", zap.String("bldgAbbr", bldgAbbr))
		query.Limit = 30 //Todo: get a definite answer on the limit
		query.Selector.ID.Regex = bldgAbbr
	} else {
		log.P.Info("getting all devices")
		query.Limit = 30 //Todo: get a definite answer on the limit
		query.Selector.ID.GT = "\x00"
	}

	var resp models.DeviceResponse
	err := couch.DBSearch(url, "POST", &query, &resp)
	if err != nil {
		log.P.Error("failed to search for devices in database")
		return nil, fmt.Errorf("Failed to find devices")
	}

	var devices []models.Device
	if resp.Docs == nil {
		log.P.Info("no devices resulted from query")
		return nil, fmt.Errorf("No devices exist under the provided search criteria")
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
	log.P.Info("searching devices by device id", zap.String("id", deviceID))
	url := fmt.Sprintf("%s/devices/%s", os.Getenv("DB_ADDRESS"), deviceID)
	var resp models.DeviceDB

	err := couch.DBSearch(url, "GET", nil, &resp)
	if err != nil {
		log.P.Error("failed to search for device in database")
		return nil, fmt.Errorf("Failed to find device with id: %s", deviceID)
	}

	s := strings.Split(resp.ID, "-")
	device := &models.Device{
		DeviceID:   resp.ID,
		DeviceName: resp.Name,
		DeviceType: resp.Type.ID,
		BldgAbbr:   s[0],
		RoomNum:    s[1],
	}
	return device, nil
}
