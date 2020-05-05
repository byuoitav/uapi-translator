package services

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/byuoitav/uapi-translator/db"
	"github.com/byuoitav/uapi-translator/log"
	"github.com/byuoitav/uapi-translator/models"
	"go.uber.org/zap"
)

//For each preset there is a master volume device
//Each mic is its own device
//Multiple outputs in one preset
//Find audioDevices in preset - take average volume returned from av api for those displays

func (s *Service) GetAudioOutputs(roomNum, bldgAbbr, devType string) ([]models.AudioOutput, error) {
	url := fmt.Sprintf("%s/ui-configuration/_find", os.Getenv("DB_ADDRESS"))
	var query models.UIConfigQuery

	if roomNum != "" && bldgAbbr != "" {
		log.Log.Info("searching audio outputs by room number and building abbreviation", zap.String("roomNum", roomNum), zap.String("bldgAbbr", bldgAbbr))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-%s$", bldgAbbr, roomNum)
	} else if roomNum != "" {
		log.Log.Info("searching audio outputs by room number", zap.String("roomNum", roomNum))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("-%s$", roomNum)
	} else if bldgAbbr != "" {
		log.Log.Info("searching audio outputs by building abbreviation", zap.String("bldgAbbr", bldgAbbr))
		query.Limit = 30
		query.Selector.ID.Regex = fmt.Sprintf("%s-", bldgAbbr)
	} else {
		log.Log.Info("getting all audio outputs")
		query.Limit = 30
		query.Selector.ID.GT = "\x00"
	}

	var resp models.AudioOutputResponse
	err := db.DBSearch(url, "POST", &query, &resp)
	if err != nil {
		log.Log.Error("failed to search for audio outputs in database")
		return nil, err
	}

	var audioOutputs []models.AudioOutput
	for _, rm := range resp.Docs {
		parts := strings.Split(rm.ID, "-")
		for i, p := range rm.Presets {

			if len(p.AudioDevices) > 0 {
				//add a master volume
				master := models.AudioOutput{
					OutputID:   fmt.Sprintf("%s-MasterAudio%d", rm.ID, (i + 1)),
					RoomNum:    parts[1],
					BldgAbbr:   parts[0],
					DeviceType: "MasterAudio",
				}
				audioOutputs = append(audioOutputs, master)
			}

			if p.IndependentAudioDevices != nil && len(p.IndependentAudioDevices) > 0 {
				for _, iad := range p.IndependentAudioDevices {
					//add the device
					deviceID := fmt.Sprintf("%s-%s", rm.ID, iad)
					device := models.AudioOutput{
						OutputID:   deviceID,
						RoomNum:    parts[1],
						BldgAbbr:   parts[0],
						DeviceType: s.getDeviceType(deviceID),
					}
					audioOutputs = append(audioOutputs, device)
				}
			}

		}
	}

	return audioOutputs, nil
}

func (s *Service) getDeviceType(devID string) string {
	device, err := s.GetDeviceByID(devID)
	if err != nil {
		return ""
	}
	return device.DeviceType
}

func (s *Service) GetAudioOutputByID(id string) (*models.AudioOutput, error) {
	log.Log.Info("searching audio outputs by id", zap.String("id", id))
	parts, index, err := s.parseOutputID(id)
	if err != nil {
		log.Log.Error("provided audio output id is invalid", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	_, err = s.getAudioOutputsFromDB(parts, index, id)
	if err != nil {
		return nil, err
	}

	var devType string
	if index == -1 {
		device, err := s.GetDeviceByID(id)
		if err != nil {
			return nil, err
		}
		devType = device.DeviceType
	} else {
		devType = "MasterAudio"
	}

	output := &models.AudioOutput{
		OutputID:   id,
		RoomNum:    parts[1],
		BldgAbbr:   parts[0],
		DeviceType: devType,
	}

	return output, nil
}

func (s *Service) GetAudioOutputState(id string) (*models.AudioOutputState, error) {
	// get ui config
	log.Log.Info("getting audio output state by id", zap.String("id", id))
	parts, index, err := s.parseOutputID(id)
	if err != nil {
		log.Log.Error("provided audio output id is invalid", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	config, err := s.getAudioOutputsFromDB(parts, index, id)
	if err != nil {
		return nil, err
	}

	//Get room state from av-api
	url := fmt.Sprintf("%s/buildings/%s/rooms/%s", os.Getenv("AV_API_URL"), parts[0], parts[1])

	var room models.RoomState
	err = db.GetState(url, "GET", &room)
	if err != nil {
		log.Log.Error("failed to find audio output state in database")
		return nil, err
	}

	if index > -1 {
		//Compare to audio devices in preset
		var volume int
		numDevices := 0
		muted := false
		for _, dev := range config.Presets[index-1].AudioDevices {
			i := s.findAudioIndex(dev, room.AudioDevices)
			if i > -1 {
				numDevices++
				volume += room.AudioDevices[i].Volume
				if room.AudioDevices[i].Muted {
					muted = true
				}
			}
		}
		//Take average of volumes
		if numDevices > 0 {
			volume /= numDevices
		}
		return &models.AudioOutputState{
			Volume: volume,
			Muted:  muted,
		}, nil
	} else {

		//Check if the id is found in the independent audio devices
		for _, p := range config.Presets {
			for _, dev := range p.IndependentAudioDevices {
				if dev == parts[2] {
					i := s.findAudioIndex(dev, room.AudioDevices)
					if i > -1 {
						return &models.AudioOutputState{
							Volume: room.AudioDevices[i].Volume,
							Muted:  room.AudioDevices[i].Muted,
						}, nil
					}
				}
			}
		}
	}

	log.Log.Infof("no state found for audio output device: %s", id)
	return nil, fmt.Errorf("no state found for audio output device: %s", id)
}

func (s *Service) findAudioIndex(name string, devices []models.StateAudioDevice) int {
	for i, dev := range devices {
		if name == dev.Name {
			return i
		}
	}
	return -1
}

func (s *Service) parseOutputID(id string) ([]string, int, error) {
	log.Log.Info("parsing audio output id", zap.String("id", id))
	parts := strings.Split(id, "-")

	if strings.Contains(parts[2], "MasterAudio") {
		index, err := strconv.Atoi(strings.Trim(parts[2], "MasterAudio"))
		if err != nil {
			return nil, 0, fmt.Errorf("Invalid audio output id")
		}

		if index < 1 {
			return nil, 0, fmt.Errorf("Invalid audio output id")
		}

		return parts, index, nil
	}
	return parts, -1, nil
}

func (s *Service) getAudioOutputsFromDB(parsedID []string, index int, id string) (*models.AudioOutputDB, error) {
	url := fmt.Sprintf("%s/ui-configuration/%s", os.Getenv("DB_ADDRESS"), fmt.Sprintf("%s-%s", parsedID[0], parsedID[1]))

	var resp models.AudioOutputDB
	err := db.DBSearch(url, "GET", nil, &resp)
	if err != nil {
		log.Log.Error("failed to find audio output config in database")
		return nil, err
	}

	if index > len(resp.Presets) {
		return nil, fmt.Errorf("Audio Output: %s does not exist", id)
	}

	return &resp, err
}
