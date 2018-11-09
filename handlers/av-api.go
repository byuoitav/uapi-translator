package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/byuoitav/av-api/base"
	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/uapi-translator/helpers"
	"github.com/fatih/color"
)

var client http.Client

func init() {
	client = http.Client{}
}

// AVSetState executes a request against the AV-API to set the state of a room.
func AVSetState(roomID string, reqBody base.PublicRoom) (toReturn base.PublicRoom, ne *nerr.E) {
	log.L.Debugf("OK! Let's call the API to set the state of %s!", roomID)

	// separate out the building and room IDs
	split := strings.Split(roomID, "-")
	building := split[0]
	room := split[1]

	// build the URL to hit the AV-API
	url := fmt.Sprintf("http://%s/buildings/%s/rooms/%s", os.Getenv("AV_API_ADDRESS"), building, room)

	// marshal the request body
	state, err := json.Marshal(reqBody)
	if err != nil {
		ne.Addf("failed to marshal the AV-API request body : %s", err.Error())
		return toReturn, ne
	}

	// create the request
	req, err := http.NewRequest("PUT", url, bytes.NewReader(state))
	if err != nil {
		ne.Addf("failed to make the request to send to the AV-API : %s", err.Error())
		return toReturn, ne
	}

	log.L.Debugf("Skittles is sending a request to %s!", url)

	// execute the request
	resp, err := client.Do(req)
	if err != nil {
		ne.Addf("failed to execute request against the AV-API : %s", err.Error())
		return toReturn, ne
	}

	// read the response
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ne.Addf("failed to read the response from the AV-API : %s", err.Error())
		return toReturn, ne
	}

	defer resp.Body.Close()

	// unmarshal the response
	err = json.Unmarshal(b, &toReturn)
	if err != nil {
		ne.Addf("failed to unmarshal the response from the AV-API : %s", err.Error())
		return toReturn, ne
	}

	log.L.Debug(color.HiCyanString("Yay! Skittles got a response from the URL %s!", url))

	return toReturn, ne
}

// AVGetState executes a request against the database to build the state of a room.
func AVGetState(roomID string) (toReturn base.PublicRoom, ne *nerr.E) {
	log.L.Debugf("Alrighty, X marks the Database! Let's get the state of %s!", roomID)

	// get the device states from the database
	deviceStates, err := db.GetDB().GetDeviceStatesByRoom(roomID)
	if err != nil {
		ne.Addf("failed to get device states for room %s from the database : %s", roomID, err.Error())
		return toReturn, ne
	}

	// get the devices from the database
	devices, err := db.GetDB().GetDevicesByRoom(roomID)
	if err != nil {
		ne.Addf("failed to get devices for room %s from the database : %s", roomID, err.Error())
		return toReturn, ne
	}

	log.L.Debugf("Skittles is reading through the device states in %s!", roomID)

	// iterate through the device states and build Displays and AudioDevices as necessary
	for _, state := range deviceStates {
		for _, device := range devices {
			if device.ID == state.DeviceID {
				if device.HasRole("VideoOut") {
					// make the display object for the room state
					display, ner := helpers.StateToDisplay(state)
					if ner != nil {
						ne.Add(ner.String())
						continue
					}

					// add the display object to the list in the room state
					toReturn.Displays = append(toReturn.Displays, display)
				}
				if device.HasRole("AudioOut") {
					// make the audio device object for the room state
					audio, ner := helpers.StateToAudioDevice(state)
					if ner != nil {
						ne.Add(ner.String())
						continue
					}

					// add the audio device to the list in the room state
					toReturn.AudioDevices = append(toReturn.AudioDevices, audio)
				}
			}
		}
	}

	split := strings.Split(roomID, "-")
	toReturn.Building = split[0]
	toReturn.Room = split[1]

	log.L.Debugf("Yay! Skittles finished reading through device states for %s!", roomID)
	return toReturn, ne
}

// AVGetConfig executes a request against the AV-API to get the configuration of a room.
func AVGetConfig(roomID string) (ne *nerr.E) {
	return ne
}
