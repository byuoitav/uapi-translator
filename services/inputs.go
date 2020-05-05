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

func (s *Service) GetInputs(roomNum, bldgAbbr string) ([]models.Input, error) {
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
		parts := strings.Split(rm.ID, "-")
		for _, in := range rm.InputConfiguration {
			deviceID := fmt.Sprintf("%s-%s", rm.ID, in.Name)
			next := models.Input{
				DeviceID:   deviceID,
				RoomNum:    parts[1],
				BldgAbbr:   parts[0],
				DeviceType: s.getDeviceType(deviceID),
				Outputs:    s.getInputDisplays(in.Name, &rm),
			}
			inputs = append(inputs, next)
		}
	}

	return inputs, nil
}

func (s *Service) GetInputByID(id string) (*models.Input, error) {
	log.Log.Info("searching inputs by id", zap.String("id", id))
	parts := strings.Split(id, "-")

	device, err := s.GetDeviceByID(id)
	if err != nil {
		log.Log.Errorf("failed to find input in database", zap.Error(err))
		return nil, err
	}

	var query models.UIConfigQuery
	query.Limit = 1000
	query.Selector.ID.Regex = fmt.Sprintf("%s-%s$", parts[0], parts[1])
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
		Outputs:    s.getInputDisplays(parts[2], &resp.Docs[0]),
	}

	return input, nil
}

func (s *Service) getInputDisplays(inputID string, resp *models.InputDB) []string {
	var displays []string
	parts := strings.Split(resp.ID, "-")
	for i, p := range resp.Presets {
		for _, in := range p.Inputs {
			if inputID == in {
				displays = append(displays, fmt.Sprintf("%s-%s-Display%d", parts[0], parts[1], (i+1)))
			}
		}
	}
	return displays
}
