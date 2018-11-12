package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/byuoitav/common/jsonhttp"
	"github.com/byuoitav/common/structs"

	"github.com/byuoitav/av-api/base"
	"github.com/byuoitav/av-api/inputgraph"
	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/v2/auth"
	"github.com/byuoitav/uapi-translator/helpers"
	"github.com/fatih/color"
)

var client http.Client

func init() {
	client = http.Client{}
}

// AVSetState executes a request against the AV-API to set the state of a room.
func AVSetState(roomID string, reqBody base.PublicRoom) (base.PublicRoom, *nerr.E) {
	log.L.Debugf("OK! Let's call the API to set the state of %s!", roomID)

	var toReturn base.PublicRoom

	// separate out the building and room IDs
	split := strings.Split(roomID, "-")
	building := split[0]
	room := split[1]

	// build the URL to hit the AV-API
	url := fmt.Sprintf("http://%s/buildings/%s/rooms/%s", os.Getenv("AV_API_ADDRESS"), building, room)

	// create the request
	req, err := jsonhttp.CreateRequest("PUT", url, reqBody, nil)
	if err != nil {
		return toReturn, nerr.Translate(err).Add("failed to make the request to send to the AV-API")
	}

	auth.AddAuthToRequest(req)

	log.L.Debugf("GoGo is sending a request to %s!", url)

	// execute the request
	resp, err := client.Do(req)
	if err != nil {
		return toReturn, nerr.Translate(err).Add("failed to execute request against the AV-API")
	}

	// read the response
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return toReturn, nerr.Translate(err).Add("failed to read the response from the AV-API")
	}

	defer resp.Body.Close()

	// unmarshal the response
	err = json.Unmarshal(b, &toReturn)
	if err != nil {
		return toReturn, nerr.Translate(err).Add("failed to unmarshal the response from the AV-API")
	}

	log.L.Debug(color.HiCyanString("Yay! GoGo got a response from the URL %s!", url))

	return toReturn, nil
}

// AVGetState executes a request against the database to build the state of a room.
func AVGetState(roomID string) (base.PublicRoom, *nerr.E) {
	log.L.Debugf("Alrighty, X marks the Database! Let's get the state of %s!", roomID)

	var toReturn base.PublicRoom

	// get the device states from the database
	deviceStates, err := db.GetDB().GetDeviceStatesByRoom(roomID)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get device states for room %s from the database", roomID)
	}

	// get the devices from the database
	devices, err := db.GetDB().GetDevicesByRoom(roomID)
	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get devices for room %s from the database : %s", roomID)
	}

	log.L.Debugf("GoGo is reading through the device states in %s!", roomID)

	// iterate through the device states and build Displays and AudioDevices as necessary
	for _, state := range deviceStates {
		for _, device := range devices {
			if device.ID == state.DeviceID {
				if device.HasRole("VideoOut") {
					// make the display object for the room state
					display, ner := helpers.StateToDisplay(state)
					if ner != nil {
						continue
					}

					// add the display object to the list in the room state
					toReturn.Displays = append(toReturn.Displays, display)
				}
				if device.HasRole("AudioOut") {
					// make the audio device object for the room state
					audio, ner := helpers.StateToAudioDevice(state)
					if ner != nil {
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
	toReturn.Room = roomID

	log.L.Debugf("Yay! GoGo finished reading through device states for %s!", roomID)
	return toReturn, nil
}

// AVGetConfig executes a request against the AV-API to get the configuration of a room.
func AVGetConfig(roomID string) (inputgraph.ReachableRoomConfig, *nerr.E) {
	log.L.Debugf("TODO: PIRATE DEBUG STATEMENT: %s", roomID)

	var toReturn inputgraph.ReachableRoomConfig

	//Get the room config from the Database
	building := strings.Split(roomID, "-")[0]
	roomName := strings.Split(roomID, "-")[1]
	room, err := db.GetDB().GetRoom(fmt.Sprintf("%s-%s", building, roomName))

	if err != nil {
		return toReturn, nerr.Translate(err).Addf("failed to get configuration for room %s from the database", roomID)
	}

	//We need the video device reachability stuff so we can filter by that
	toReturn, ne := inputgraph.GetVideoDeviceReachability(room)
	if err != nil {
		return toReturn, ne
	}

	//TODO Filter out devices that don't have AudioOut, AudioIn, or aren't in the InputReachability Map
	//I think this should filter out non-audio(in/out) stuff. Now we need to filter out the stuff outside of the input reachability map
	//And get rid of the other stuff that's marked out
	var devices []structs.Device
	for _, device := range toReturn.Devices {
		for _, role := range device.Roles {
			if role.ID == "AudioIn" || role.ID == "AudioOut" {
				devices = append(devices, device)
			} else {
				//Just to show what is not getting added, can be removed
				log.L.Infof("%s", device.ID)
			}
		}
	}

	//Set the devices to the filtered devices
	toReturn.Devices = devices

	log.L.Debugf("Avast ye! Here lies yer configuration for %s!", roomID)
	return toReturn, nil
}
