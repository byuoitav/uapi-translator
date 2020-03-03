package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/byuoitav/uapi-translator/couch"
	"github.com/byuoitav/uapi-translator/models"
)

//From ui config
//Get presets
//Create name for each preset group

func GetDisplays(roomNum, bldgAbbr string) ([]models.Display, error) {
	url := fmt.Sprintf("%s/ui-configuration/_find", os.Getenv("DB_ADDRESS"))
	var query models.DisplayQuery

	if roomNum != "" && bldgAbbr != "" {
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-%s$", bldgAbbr, roomNum)
	} else if roomNum != "" {
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("-%s$", roomNum)
	} else if bldgAbbr != "" {
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-", bldgAbbr)
	} else {
		query.Limit = 30
		query.Selector.ID.GT = "\x00"
	}
	//post query

	var resp models.DisplayResponse
	err := couch.DBSearch(url, "POST", &query, &resp)
	if err != nil {
		return nil, err
	}

	var displays []models.Display
	if resp.Docs == nil {
		return nil, fmt.Errorf("No displays")
	}

	for _, rm := range resp.Docs {
		for j := range rm.Presets {
			s := strings.Split(rm.ID, "-")
			next := models.Display{
				DisplayID: fmt.Sprintf("%s-Display%d", rm.ID, (j + 1)),
				RoomNum:   s[1],
				BldgAbbr:  s[0],
			}
			displays = append(displays, next)
		}
	}

	return displays, nil
}
