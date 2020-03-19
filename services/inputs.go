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
				DeviceID: deviceID,
				RoomNum: s[1],
				BldgAbbr: s[0],
				DeviceType: getDeviceType(deviceID),
				Outputs: nil,
			}
			inputs = append(inputs, next)
		}
	}

	return inputs, nil
}
