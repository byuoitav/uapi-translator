package services

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/byuoitav/scheduler/log"
	"github.com/byuoitav/uapi-translator/db"
	"github.com/byuoitav/uapi-translator/models"
)

func GetDisplays(roomNum, bldgAbbr string) ([]models.Display, error) {
	url := fmt.Sprintf("%s/ui-configuration/_find", os.Getenv("DB_ADDRESS"))
	var query models.DisplayQuery

	if roomNum != "" && bldgAbbr != "" {
		log.P.Info("searching displays by room number and building abbreviation", zap.String("roomNum", roomNum), zap.String("bldgAbbr", bldgAbbr))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-%s$", bldgAbbr, roomNum)
	} else if roomNum != "" {
		log.P.Info("searching displays by room number", zap.String("roomNum", roomNum))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("-%s$", roomNum)
	} else if bldgAbbr != "" {
		log.P.Info("searching displays by building abbreviation", zap.String("bldgAbbr", bldgAbbr))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-", bldgAbbr)
	} else {
		log.P.Info("getting all displays")
		query.Limit = 30
		query.Selector.ID.GT = "\x00"
	}

	var resp models.DisplayResponse
	err := db.DBSearch(url, "POST", &query, &resp)
	if err != nil {
		log.P.Error("failed to search for displays in database")
		return nil, err
	}

	var displays []models.Display
	if resp.Docs == nil {
		log.P.Info("no displays resulted from query")
		return nil, fmt.Errorf("No displays exist under the provided search criteria")
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

func GetDisplayByID(dispID string) (*models.Display, error) {
	log.P.Info("searching displays by display id", zap.String("id", dispID))
	s, index, err := parseDisplayID(dispID)
	if err != nil {
		log.P.Error("provided display id is invalid", zap.String("id", dispID), zap.Error(err))
		return nil, err
	}

	_, err = getDisplaysFromDB(s, index, dispID)
	if err != nil {
		return nil, err
	}

	display := &models.Display{
		DisplayID: dispID,
		RoomNum:   s[1],
		BldgAbbr:  s[0],
	}
	return display, nil
}

func GetDisplayConfig(dispID string) (*models.DisplayConfig, error) {
	log.P.Info("searching for display config", zap.String("id", dispID))
	s, index, err := parseDisplayID(dispID)
	if err != nil {
		log.P.Error("provided display id is invalid", zap.String("id", dispID), zap.Error(err))
		return nil, err
	}

	displays, err := getDisplaysFromDB(s, index, dispID)
	if err != nil {
		return nil, err
	}

	var devices []string
	for _, dev := range displays.Presets[index-1].Displays {
		devices = append(devices, fmt.Sprintf("%s-%s-%s", s[0], s[1], dev))
	}

	var inputs []string
	for _, in := range displays.Presets[index-1].Inputs {
		inputs = append(inputs, fmt.Sprintf("%s-%s-%s", s[0], s[1], in))
	}

	config := &models.DisplayConfig{
		Devices: devices,
		Inputs:  inputs,
	}

	return config, nil
}

func GetDisplayState(dispID string) (*models.DisplayState, error) {
	log.P.Info("searching for display state", zap.String("id", dispID))
	s, index, err := parseDisplayID(dispID)
	if err != nil {
		log.P.Error("provided display id is invalid", zap.String("id", dispID), zap.Error(err))
		return nil, err
	}

	//send request to av api
	url := fmt.Sprintf("%s/buildings/%s/rooms/%s", os.Getenv("AV_API_URL"), s[0], s[1])

	var resp models.RoomState
	err = db.GetState(url, "GET", &resp)
	if err != nil {
		log.P.Error("failed to find display state in database")
		return nil, err
	}

	displays, err := getDisplaysFromDB(s, index, dispID)
	if err != nil {
		return nil, err
	}

	//compare displays within the same preset
	var firstDisplay models.StateDisplays
	for i, p := range displays.Presets[index-1].Displays {
		if i == 0 {
			firstDisplay = 
		} else {

		}
	}

	state := &models.DisplayState{
		Powered: true,
		Blanked: true,
		Input:   "test",
	}
	return state, nil
}

func parseDisplayID(id string) ([]string, int, error) {
	log.P.Info("parsing display id", zap.String("id", id))
	s := strings.Split(id, "-")

	if !strings.Contains(s[2], "Display") {
		return nil, 0, fmt.Errorf("Invalid display id")
	}

	index, err := strconv.Atoi(strings.Trim(s[2], "Display"))
	if err != nil {
		return nil, 0, fmt.Errorf("Invalid display id")
	}

	return s, index, nil
}

func findDisplayIndex(id string, obj *models.RoomState) int {
	for index, disp := range obj. {
		if id == disp {
			return index
		}
	}
	return -1
}

func getDisplaysFromDB(parsedID []string, index int, dispID string) (*models.DisplayDB, error) {
	url := fmt.Sprintf("%s/ui-configuration/%s", os.Getenv("DB_ADDRESS"), fmt.Sprintf("%s-%s", parsedID[0], parsedID[1]))

	var resp models.DisplayDB
	err := db.DBSearch(url, "GET", nil, &resp)
	if err != nil {
		log.P.Error("failed to find display config in database")
		return nil, err
	}

	if index > len(resp.Presets) {
		return nil, fmt.Errorf("Display: %s does not exist", dispID)
	}

	return &resp, err
}
