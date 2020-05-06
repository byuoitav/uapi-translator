package services

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/byuoitav/uapi-translator/db"
	"github.com/byuoitav/uapi-translator/log"
	"github.com/byuoitav/uapi-translator/models"
)

func (s *Service) GetRooms(roomNum, bldgAbbr string) ([]models.Room, error) {
	url := fmt.Sprintf("%s/rooms/_find", os.Getenv("DB_ADDRESS"))
	var query models.RoomQuery

	if roomNum != "" && bldgAbbr != "" {
		log.Log.Info("searching rooms by room number and building abbreviation", zap.String("roomNum", roomNum), zap.String("bldgAbbr", bldgAbbr))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("%s-%s$", bldgAbbr, roomNum)
	} else if roomNum != "" {
		log.Log.Info("searching rooms by room number", zap.String("roomNum", roomNum))
		query.Limit = 1000
		query.Selector.ID.Regex = fmt.Sprintf("-%s$", roomNum)
	} else if bldgAbbr != "" {
		log.Log.Info("searching rooms by building abbreviation", zap.String("bldgAbbr", bldgAbbr))
		query.Limit = 30 //Todo: get a definite answer on the limit
		query.Selector.ID.Regex = bldgAbbr
	} else {
		log.Log.Info("getting all rooms")
		query.Limit = 30 //Todo: get a definite answer on the limit
		query.Selector.ID.GT = "\x00"
	}

	var resp db.RoomResponse
	err := db.DBSearch(url, "POST", &query, &resp)
	if err != nil {
		log.Log.Error("failed to search for rooms in database", zap.Error(err))
		return nil, err
	}

	var rooms []models.Room
	if resp.Docs == nil {
		log.Log.Info("no rooms resulted from query")
		return nil, fmt.Errorf("No rooms exist under the provided search criteria")
	}
	for _, rm := range resp.Docs {
		roomParts := strings.Split(rm.ID, "-")
		resources, err := s.GetRoomResources(rm.ID)
		if err != nil {
			return nil, fmt.Errorf("services/GetRooms get room resources: %w", err)
		}
		next := models.Room{
			RoomID:      rm.ID,
			RoomNum:     roomParts[1],
			BldgAbbr:    roomParts[0],
			Description: rm.Tags["description"],
			Resources:   resources,
		}
		rooms = append(rooms, next)
	}
	return rooms, nil
}

func (s *Service) GetRoomDevices(roomID string) (*models.RoomDevices, error) {
	// Check if room exists
	roomParts := strings.Split(roomID, "-")
	_, err := s.GetRooms(roomParts[1], roomParts[0])
	if err != nil {
		return nil, fmt.Errorf("No rooms exist with the id: %s", roomID)
	}

	var devices models.RoomDevices
	displays, err := s.GetDisplays(roomParts[1], roomParts[0])
	if err == nil {
		for _, disp := range displays {
			devices.Displays = append(devices.Displays, disp.DisplayID)
		}
	}

	audioOutputs, err := s.GetAudioOutputs(roomParts[1], roomParts[0], "")
	if err == nil {
		for _, out := range audioOutputs {
			devices.Outputs = append(devices.Outputs, out.OutputID)
		}
	}

	inputs, err := s.GetInputs(roomParts[1], roomParts[0])
	if err == nil {
		for _, in := range inputs {
			devices.Inputs = append(devices.Inputs, in.DeviceID)
		}
	}

	return &devices, nil
}

// GetRoomResources returns an array of the resources associated with
// the given roomID
func (s *Service) GetRoomResources(roomID string) ([]models.Resource, error) {

	devs, err := s.DB.GetDevicesByRoom(roomID)
	if err != nil {
		return nil, fmt.Errorf("services/GetRoomResources get devices: %w", err)
	}

	types := map[string]*db.DeviceType{}
	resources := map[string]models.Resource{}

	// Abstract resources from the devices
	for _, d := range devs {
		// Get description
		desc := ""
		// Check for description tag on device
		if val, ok := d.Tags["description"]; ok {
			desc = val
		} else if t, ok := types[d.TypeID]; ok {
			// If we have already pulled the type then use its description field
			desc = t.Tags["description"]
		} else {
			// If we haven't pulled the type then pull it and use its description
			t, err := s.DB.GetDeviceTypeByID(d.TypeID)
			if err != nil {
				return nil, fmt.Errorf("services/GetRoomResources get device type: %w", err)
			}
			types[t.ID] = t
			desc = t.Tags["description"]
		}

		// Skip indescribable objects
		if desc == "" {
			continue
		}

		// Check to see if we are already tracking this resource type
		if val, ok := resources[desc]; ok {
			val.Quantity += 1 // increment quantity
			// Append location if there is one
			if d.Tags["location"] != "" {
				val.Locations = append(val.Locations, d.Tags["location"])
			}

			// Update
			resources[desc] = val
		} else { // This is the first resource of it's kind in this list
			r := models.Resource{}
			r.Quantity = 1
			r.Resource = desc
			r.Locations = []string{}

			// Append location if there is one
			if d.Tags["location"] != "" {
				r.Locations = append(r.Locations, d.Tags["location"])
			}

			// Put it in the map
			resources[desc] = r
		}
	}

	r := make([]models.Resource, len(resources))
	i := 0
	for _, resource := range resources {
		r[i] = resource
		i++
	}

	return r, nil
}
