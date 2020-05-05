package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

	res, err := s.makeRequest("GET", path, nil)
	if err != nil {
		err = fmt.Errorf("db/GetRoomByID make request: %w", err)
		return nil, err
	}
	defer res.Body.Close()

	// Check for 404
	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	// Check for non 200
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error response from couch. Code: %d", res.StatusCode)
	}

	// Read the body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("db/GetRoomID reading response body: %w", err)
	}

	// Unmarshal
	room := Room{}
	err = json.Unmarshal(body, &room)
	if err != nil {
		return nil, fmt.Errorf("db/GetRoomID json unmarshal: %w", err)
	}

	return &room, nil
}
