package models

import "github.com/byuoitav/common/structs"

// Couch query object
type PrefixQuery struct {
	Selector struct {
		ID struct {
			GT    string `json:"$gt,omitempty"`
			LT    string `json:"$lt,omitempty"`
			Regex string `json:"$regex,omitempty"`
		} `json:"_id"`
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
