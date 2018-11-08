package helpers

import (
	"fmt"
	"strings"

	"github.com/byuoitav/av-api/base"
	"github.com/byuoitav/common/nerr"
	sd "github.com/byuoitav/common/state/statedefinition"
	"github.com/byuoitav/uapi-translator/structs"
)

// Field Set constants
const (
	Basic    = "basic"
	State    = "av_state"
	Config   = "av_config"
	SetState = "set_state"
	GetState = "get_state"
)

// API Type constants
const (
	ReadOnly     = "read-only"
	Modifiable   = "modifiable"
	System       = "system"
	Derived      = "derived"
	Unauthorized = "unauthorized"
	Related      = "related"
)

// NotFound is a response code
var NotFound = 404

// Okay is a response code
var Okay = 201

var apiURL = "https://api.byu.edu/byuapi/av_rooms"

var available = []string{"basic", "av_state", "av_config"}

var def = []string{"basic"}

// UAPItoAV takes the Resource structure of the UAPI and translates it into the AV-API PublicRoom structure.
func UAPItoAV(resource structs.Resource) (pubRoom base.PublicRoom, ne *nerr.E) {
	return pubRoom, ne
}

// AVtoUAPI takes the AV-API PublicRoom structure and translates it into the UAPI Resource structure.
func AVtoUAPI(pubRoom base.PublicRoom, fieldSets ...string) (resource structs.Resource, ne *nerr.E) {
	// combine building and room to make a roomID
	roomID := fmt.Sprintf("%s-%s", pubRoom.Building, pubRoom.Room)

	// create top-level resource Metadata
	resource.Metadata.FieldSetsAvailable = available
	resource.Metadata.FieldSetsReturned = fieldSets
	resource.Metadata.FieldSetsDefault = def
	resource.Metadata.ValidationResponse = structs.ValidationResponse{}

	resource.Basic = structs.SubResource{}
	resource.State = structs.SubResource{}
	resource.Config = structs.SubResource{}

	// create SubResources for each field set
	for _, fs := range fieldSets {
		switch fs {
		case Basic:
			resource.Basic = generateBasicSubResource(roomID, pubRoom)
			break
		case State:
			resource.State = generateStateSubResource(roomID, pubRoom)
			break
		case Config:
			resource.Config = generateConfigSubResource(roomID, pubRoom)
			break
		default:
			break
		}
	}

	return resource, ne
}

func generateBasicSubResource(roomID string, pubRoom base.PublicRoom) (sub structs.SubResource) {
	// set the metadata first, because if the values are not correct then we should return only the metadata
	if len(pubRoom.Building) == 0 || len(pubRoom.Room) == 0 {
		sub.Metadata.ValidationResponse = structs.ValidationResponse{
			Code:    &NotFound,
			Message: "resource not found, invalid room information",
		}
		return sub
	}

	// create the metadata for the subresource
	sub.Metadata.ValidationResponse = structs.ValidationResponse{
		Code:    &Okay,
		Message: "success",
	}

	// create links for the subresource
	var link structs.Link

	link.Rel = "self"
	link.Href = fmt.Sprintf("%s/%s", apiURL, roomID)
	link.Method = "GET"

	sub.Links = make(map[string]structs.Link)

	sub.Links["basics__info"] = link

	// set the Building and Room fields
	sub.Building = structs.Property{
		Type:        ReadOnly,
		Key:         true,
		Value:       pubRoom.Building,
		ValueArray:  nil,
		Object:      nil,
		ObjectArray: nil,
	}

	sub.Room = structs.Property{
		Type:        ReadOnly,
		Key:         true,
		Value:       roomID,
		ValueArray:  nil,
		Object:      nil,
		ObjectArray: nil,
	}

	return sub
}

func generateStateSubResource(roomID string, pubRoom base.PublicRoom) (sub structs.SubResource) {
	// set the metadata first, because if the values are not correct then we should return only the metadata
	if len(pubRoom.Building) == 0 || len(pubRoom.Room) == 0 {
		sub.Metadata.ValidationResponse = structs.ValidationResponse{
			Code:    &NotFound,
			Message: "resource not found, invalid room information",
		}
		return sub
	}

	if len(pubRoom.Displays) == 0 || len(pubRoom.AudioDevices) == 0 {
		sub.Metadata.ValidationResponse = structs.ValidationResponse{
			Code:    &NotFound,
			Message: "unable to get state of the room",
		}
		return sub
	}

	// create the metadata for the subresource
	sub.Metadata.ValidationResponse = structs.ValidationResponse{
		Code:    &Okay,
		Message: "success",
	}

	// create links for the subresource
	var get structs.Link
	var set structs.Link

	get.Rel = State
	set.Rel = State

	get.Href = fmt.Sprintf("%s/%s/%s", apiURL, roomID, State)
	set.Href = fmt.Sprintf("%s/%s/%s", apiURL, roomID, State)

	get.Method = "GET"
	set.Method = "PUT"

	sub.Links = make(map[string]structs.Link)

	sub.Links[GetState] = get
	sub.Links[SetState] = set

	// set the properties on the subresource
	sub.Building = structs.Property{
		Type:        ReadOnly,
		Key:         true,
		Value:       pubRoom.Building,
		ValueArray:  nil,
		Object:      nil,
		ObjectArray: nil,
	}

	sub.Room = structs.Property{
		Type:        ReadOnly,
		Key:         true,
		Value:       roomID,
		ValueArray:  nil,
		Object:      nil,
		ObjectArray: nil,
	}

	// create a list of Properties based on the displays.
	var dispList = make(map[string]interface{})

	for _, display := range pubRoom.Displays {
		p := structs.Property{
			Type: Modifiable,
			ObjectArray: map[string]interface{}{
				"name": structs.Property{
					Type:        ReadOnly,
					Key:         true,
					Value:       display.Name,
					ValueArray:  nil,
					Object:      nil,
					ObjectArray: nil,
				},
				"power": structs.Property{
					Type:        Modifiable,
					Value:       display.Power,
					ValueArray:  nil,
					Object:      nil,
					ObjectArray: nil,
				},
				"input": structs.Property{
					Type:        Modifiable,
					Value:       display.Input,
					ValueArray:  nil,
					Object:      nil,
					ObjectArray: nil,
				},
				"blanked": structs.Property{
					Type:        Modifiable,
					Value:       fmt.Sprintf("%b", display.Blanked),
					ValueArray:  nil,
					Object:      nil,
					ObjectArray: nil,
				},
			},
			ValueArray: nil,
			Object:     nil,
		}

		dispList[display.Name] = p
	}

	// create a list of Properties based on the Audio Devices
	var audioList = make(map[string]interface{})

	for _, audio := range pubRoom.AudioDevices {
		a := structs.Property{
			Type: Modifiable,
			ObjectArray: map[string]interface{}{
				"name": structs.Property{
					Type:  ReadOnly,
					Key:   true,
					Value: audio.Name,
				},
				"power": structs.Property{
					Type:  Modifiable,
					Value: audio.Power,
				},
				"input": structs.Property{
					Type:  Modifiable,
					Value: audio.Input,
				},
				"muted": structs.Property{
					Type:  Modifiable,
					Value: fmt.Sprintf("%b", audio.Muted),
				},
				"volume": structs.Property{
					Type:  Modifiable,
					Value: fmt.Sprintf("%d", audio.Volume),
				},
			},
			ValueArray: nil,
			Object:     nil,
		}

		audioList[audio.Name] = a
	}

	sub.Displays = structs.Property{
		Type:        Modifiable,
		ObjectArray: dispList,
		ValueArray:  nil,
		Object:      nil,
	}

	sub.AudioDevices = structs.Property{
		Type:        Modifiable,
		ObjectArray: audioList,
		ValueArray:  nil,
		Object:      nil,
	}

	return sub
}

func generateConfigSubResource(roomID string, pubRoom base.PublicRoom) (sub structs.SubResource) {
	// create links for the subresource
	var con structs.Link

	con.Rel = "self"
	con.Href = fmt.Sprintf("%s/%s/%s", apiURL, roomID, Config)
	con.Method = "GET"

	sub.Links = make(map[string]structs.Link)

	sub.Links[Config] = con
	return sub
}

// StateToDisplay takes in a DeviceState object and translates it into a base.Display
func StateToDisplay(state sd.StaticDevice) (display base.Display, ne *nerr.E) {
	// set the display Name if the state DeviceID is not empty
	if len(state.DeviceID) > 0 {
		display.Name = strings.Split(state.DeviceID, "-")[2]
	}

	// set the display Power if the state Power is not empty
	if len(state.Power) > 0 {
		display.Power = state.Power
	}

	// set the display Input if the state Input is not empty
	if len(state.Input) > 0 {
		display.Input = state.Input
	}

	// set the display Blanked if the state Blanked is not empty
	if state.Blanked != nil {
		display.Blanked = state.Blanked
	}

	return display, ne
}

// StateToAudioDevice takes in a DeviceState object and translates it into a base.AudioDevice
func StateToAudioDevice(state sd.StaticDevice) (audio base.AudioDevice, ne *nerr.E) {
	// set the audio Name if the state DeviceID is not empty
	if len(state.DeviceID) > 0 {
		audio.Name = strings.Split(state.DeviceID, "-")[2]
	}

	// set the audio Power if the state Power is not empty
	if len(state.Power) > 0 {
		audio.Power = state.Power
	}

	// set the audio Input if the state Input is not empty
	if len(state.Input) > 0 {
		audio.Input = state.Input
	}

	// set the audio Muted if the state Muted is not empty
	if state.Muted != nil {
		audio.Muted = state.Muted
	}

	// set the audio Volume if the state Volume is not empty
	if state.Volume != nil {
		audio.Volume = state.Volume
	}

	return audio, ne
}
