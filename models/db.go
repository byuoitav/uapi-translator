package models

import "github.com/byuoitav/common/structs"

type CouchSearch struct {
	GT    string `json:"$gt,omitempty"`
	LT    string `json:"$lt,omitempty"`
	Regex string `json:"$regex,omitempty"`
}

type RoomQuery struct {
	Selector struct {
		ID CouchSearch `json:"_id"`
	} `json:"selector"`
	Limit int `json:"limit"`
}

type DeviceTypeQuery struct {
	ID *CouchSearch `json:"_id,omitempty"`
}

type DeviceQuery struct {
	Selector struct {
		ID      CouchSearch      `json:"_id"`
		DevType *DeviceTypeQuery `json:"type,omitempty"`
	} `json:"selector"`
	Limit int `json:"limit"`
}

type UIConfigQuery struct {
	Selector struct {
		ID CouchSearch `json:"_id"`
	} `json:"selector"`
	Limit int `json:"limit"`
}

// Rooms
type RoomResponse struct {
	Docs     []RoomDB `json:"docs"`
	Bookmark string   `json:"bookmark"`
	Warning  string   `json:"warning"`
}

type RoomDB struct {
	Rev string `json:"_rev,omitempty"`
	*structs.Room
}

type DeviceResponse struct {
	Docs     []DeviceDB `json:"docs"`
	Bookmark string     `json:"bookmark"`
	Warning  string     `json:"warning"`
}

type DeviceDB struct {
	Rev string `json:"_rev,omitempty"`
	*structs.Device
}

type DisplayResponse struct {
	Docs     []DisplayDB `json:"docs"`
	Bookmark string      `json:"bookmark"`
	Warning  string      `json:"warning"`
}

type DisplayDB struct {
	Rev string `json:"_rev,omitempty"`
	*structs.UIConfig
}

type AudioOutputResponse struct {
	Docs     []AudioOutputDB `json:"docs"`
	Bookmark string          `json:"bookmark"`
	Warning  string          `json:"warning"`
}

type AudioOutputDB struct {
	Rev string `json:"_rev,omitempty"`
	*structs.UIConfig
}

type RoomState struct {
	Displays     []StateDisplay     `json:"displays,omitempty"`
	AudioDevices []StateAudioDevice `json:"audioDevices,omitempty"`
}

type StateDisplay struct {
	Name    string `json:"name,omitempty"`
	Power   string `json:"power,omitempty"`
	Input   string `json:"input,omitempty"`
	Blanked bool   `json:"blanked,omitempty"`
}

type StateAudioDevice struct {
	Name   string `json:"name,omitempty"`
	Power  string `json:"power,omitempty"`
	Input  string `json:"input,omitempty"`
	Muted  bool   `json:"muted,omitempty"`
	Volume int    `json:"volume,omitempty"`
}

type InputResponse struct {
	Docs     []InputDB `json:"docs"`
	Bookmark string    `json:"bookmark"`
	Warning  string    `json:"warning"`
}

type InputDB struct {
	Rev string `json:"_rev,omitempty"`
	*structs.UIConfig
}
