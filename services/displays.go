package services

import (
	"fmt"
	"os"
	"strconv"
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

func GetDisplayByID(dispID string) (*models.Display, error) {
	s, index, err := parseDisplayID(dispID)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/ui-configuration/%s", os.Getenv("DB_ADDRESS"), fmt.Sprintf("%s-%s", s[0], s[1]))

	var resp models.DisplayDB
	err = couch.DBSearch(url, "GET", nil, &resp)
	if err != nil {
		return nil, err
	}

	if index > len(resp.Presets) {
		return nil, fmt.Errorf("Display does not exist")
	}

	display := &models.Display{
		DisplayID: dispID,
		RoomNum:   s[1],
		BldgAbbr:  s[0],
	}
	return display, nil
}

func GetDisplayConfig(dispID string) (*models.DisplayConfig, error) {
	s, index, err := parseDisplayID(dispID)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/ui-configuration/%s", os.Getenv("DB_ADDRESS"), fmt.Sprintf("%s-%s", s[0], s[1]))

	var resp models.DisplayDB
	err = couch.DBSearch(url, "GET", nil, &resp)
	if err != nil {
		return nil, err
	}

	if index > len(resp.Presets) {
		return nil, fmt.Errorf("Display does not exist")
	}

	var devices []string
	for _, dev := range resp.Presets[index-1].Displays {
		devices = append(devices, fmt.Sprintf("%s-%s-%s", s[0], s[1], dev))
	}

	var inputs []string
	for _, in := range resp.Presets[index-1].Inputs {
		inputs = append(inputs, fmt.Sprintf("%s-%s-%s", s[0], s[1], in))
	}

	config := &models.DisplayConfig{
		Devices: devices,
		Inputs:  inputs,
	}

	return config, nil
}

func parseDisplayID(id string) ([]string, int, error) {
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
