package services

import (
	"fmt"

	"github.com/byuoitav/uapi-translator/models"
)

//From ui config
//Get presets
//Create name for each preset group

func GetDisplays(roomNum, bldgAbbr string) ([]models.Display, error) {
	// url := fmt.Sprintf("%s/ui-configuration/_find", os.Getenv("DB_ADDRESS"))
	var query models.DisplayQuery

	if roomNum != "" && bldgAbbr != "" {
		roomID := fmt.Sprintf("%s-%s", bldgAbbr, roomNum)
		query.Limit = 1000
		query.Selector.ID.Regex = roomID
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
	// dbDisplays, err := requestDisplaySearch(url, "POST", &query)
	// if err != nil {
	// 	return nil, err
	// }
	//translate to models.Display
	return nil, nil
}

// func requestDisplaySearch(url, method string, query interface{}) ([]models.DisplayDB, error) {
// 	var body []byte
// 	var err error
// 	if query != nil {
// 		body, err = json.Marshal(query)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	var resp models.DisplayResponse
// 	err = couch.MakeRequest(method, url, "application/json", body, &resp)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return resp.Docs, nil
// }
