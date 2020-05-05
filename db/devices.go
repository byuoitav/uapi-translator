package db

import (
	"encoding/json"
	"fmt"
)

const _devicesPath = "devices"
const _deviceTypesPath = "device-types"

type DeviceResponse struct {
	Docs     []Device `json:"docs"`
	Bookmark string   `json:"bookmark"`
	Warning  string   `json:"warning"`
}

type Device struct {
	ID     string            `json:"_id"`
	TypeID string            `json:"typeID"`
	Tags   map[string]string `json:"tags"`
}

type DeviceType struct {
	ID   string            `json:"_id"`
	Tags map[string]string `json:"tags"`
}

// GetDevicebyID gets a device document from couch given the id
func (s *Service) GetDeviceByID(deviceID string) (*Device, error) {
	path := fmt.Sprintf("%s/%s", _devicesPath, deviceID)
	d := Device{}

	// Make request
	err := s.makeRequest("GET", path, nil, &d)
	if err != nil {
		err = fmt.Errorf("db/GetDeviceByID make request: %w", err)
		return nil, err
	}

	return &d, nil
}

// GetDevicesByRoom returns an array of devices that are a part of the given room
func (s *Service) GetDevicesByRoom(roomID string) ([]Device, error) {
	path := fmt.Sprintf("%s/_find", _devicesPath)
	r := DeviceResponse{}

	// Format query
	q := query{
		Selector: map[string]interface{}{
			"_id": search{
				Regex: fmt.Sprintf("%s-", roomID),
			},
		},
		Limit: 1000,
	}
	body, err := json.Marshal(&q)
	if err != nil {
		return nil, fmt.Errorf("db/GetDevicesByRoom query marshal: %w", err)
	}

	// Make the request
	err = s.makeRequest("POST", path, body, &r)
	if err != nil {
		return nil, fmt.Errorf("db/GetDevicesByRoom couch request: %w", err)
	}

	return r.Docs, nil
}

// GetDeviceTypeByID returns the device type document for the given id
func (s *Service) GetDeviceTypeByID(deviceTypeID string) (*DeviceType, error) {
	path := fmt.Sprintf("%s/%s", _deviceTypesPath, deviceTypeID)
	d := DeviceType{}

	// Make request
	err := s.makeRequest("GET", path, nil, &d)
	if err != nil {
		return nil, fmt.Errorf("db/GetDeviceTypeByID couch request: %w", err)
	}

	return &d, nil
}
