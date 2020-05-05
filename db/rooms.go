package db

import (
	"fmt"
)

const _roomsPath = "rooms"

// Rooms
type RoomResponse struct {
	Docs     []Room `json:"docs"`
	Bookmark string `json:"bookmark"`
	Warning  string `json:"warning"`
}

type Room struct {
	Rev  string            `json:"_rev,omitempty"`
	ID   string            `json:"_id"`
	Tags map[string]string `json:"tags"`
}

// GetRoomByID returns the Room document for the given roomID
func (s *Service) GetRoomByID(roomID string) (*Room, error) {
	path := fmt.Sprintf("%s/%s", _roomsPath, roomID)

	room := Room{}
	err := s.makeRequest("GET", path, nil, &room)
	if err != nil {
		err = fmt.Errorf("db/GetRoomByID make request: %w", err)
		return nil, err
	}

	return &room, nil
}
