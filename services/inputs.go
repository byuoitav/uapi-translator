package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/byuoitav/uapi-translator/db"
	"github.com/byuoitav/uapi-translator/log"
	"github.com/byuoitav/uapi-translator/models"
	"go.uber.org/zap"
)

func GetInputs(roomNum, bldgAbbr string) ([]models.Input, error) {
	url := fmt.Sprintf("%s/ui-configuration/_find", os.Getenv("DB_ADDRESS"))
	var query models.UIConfigQuery

	if roomNum != "" && bldgAbbr != "" {
		log.Log.Info("searching inputs by room number and building abbreviation", zap.String("roomNum", roomNum), zap.String("bldgAbbr", bldgAbbr))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-%s$", bldgAbbr, roomNum)
	} else if roomNum != "" {
		log.Log.Info("searching inputs by room number", zap.String("roomNum", roomNum))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("-%s$", roomNum)
	} else if bldgAbbr != "" {
		log.Log.Info("searching inputs by building abbreviation", zap.String("bldgAbbr", bldgAbbr))
		query.Limit = 30
		query.Selector.ID.Regex = fmt.Sprintf("%s-", bldgAbbr)
	} else {
		log.Log.Info("getting all inputs")
		query.Limit = 30
		query.Selector.ID.GT = "\x00"
	}

	var resp models.InputResponse
	err := db.DBSearch(url, "POST", &query, &resp)
	if err != nil {
		log.Log.Error("failed to search for inputs in database")
		return nil, err
	}

	var inputs []models.Input
	for _, rm := range resp.Docs {
		s := strings.Split(rm.ID, "-")
		for _, in := range rm.InputConfiguration {
			deviceID := fmt.Sprintf("%s-%s", rm.ID, in.Name)
			next := models.Input{
				DeviceID:   deviceID,
				RoomNum:    s[1],
				BldgAbbr:   s[0],
				DeviceType: getDeviceType(deviceID),
				Outputs:    getInputDisplays(in.Name, &rm),
			}
			inputs = append(inputs, next)
		}
	}

	return inputs, nil
}

func GetInputByID(id string) (*models.Input, error) {
	log.Log.Info("searching inputs by id", zap.String("id", id))
	s := strings.Split(id, "-")

	device, err := GetDeviceByID(id)
	if err != nil {
		log.Log.Errorf("failed to find input in database", zap.Error(err))
		return nil, err
	}

	var query models.UIConfigQuery
	query.Limit = 1000
	query.Selector.ID.Regex = fmt.Sprintf("%s-%s$", s[0], s[1])
	url := fmt.Sprintf("%s/ui-configuration/_find", os.Getenv("DB_ADDRESS"))

	var resp models.InputResponse
	err = db.DBSearch(url, "POST", &query, &resp)
	if err != nil {
		log.Log.Error("failed to search for input in database")
		return nil, err
	}

	input := &models.Input{
		DeviceID:   device.DeviceID,
		RoomNum:    device.RoomNum,
		BldgAbbr:   device.BldgAbbr,
		DeviceType: device.DeviceType,
		Outputs:    getInputDisplays(s[2], &resp.Docs[0]),
	}

	return input, nil
}

func getInputDisplays(inputID string, resp *models.InputDB) []string {
	var displays []string
	s := strings.Split(resp.ID, "-")
	for i, p := range resp.Presets {
		for _, in := range p.Inputs {
			if inputID == in {
				displays = append(displays, fmt.Sprintf("%s-%s-Display%d", s[0], s[1], (i + 1)))
			}
		}
	}
	return displays
}
