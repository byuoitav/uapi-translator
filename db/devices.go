package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

	// Make request
	res, err := s.makeRequest("GET", path, nil)
	if err != nil {
		err = fmt.Errorf("db/GetDeviceByID make request: %w", err)
		return nil, err
	}
	defer res.Body.Close()

	// Check if the document doesn't exist
	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	// Check for error code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("db/GetDeviceByID got non 200 from couch. Code: %d", res.StatusCode)
	}

	// Read the body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("db/GetDeviceByID read body: %w", err)
		return nil, err
	}

	// Unmarshal
	d := Device{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		err = fmt.Errorf("db/GetDeviceByID json umarshal: %w", err)
		return nil, err
	}

	return &d, nil
}

// GetDevicesByRoom returns an array of devices that are a part of the given room
func (s *Service) GetDevicesByRoom(roomID string) ([]Device, error) {
	path := fmt.Sprintf("%s/_find", _devicesPath)

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
	res, err := s.makeRequest("POST", path, body)
	if err != nil {
		return nil, fmt.Errorf("db/GetDevicesByRoom couch request: %w", err)
	}
	defer res.Body.Close()

	// Check for error code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("db/GetDevicesByRoom got non 200 from couch: Code: %d", res.StatusCode)
	}

	// Read body
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("db/GetDevicesByRoom read body: %w", err)
	}

	// JSON unmarshal
	r := DeviceResponse{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, fmt.Errorf("db/GetDevicesByRoom json unmarshal: %w", err)
	}

	return r.Docs, nil
}

// GetDeviceTypeByID returns the device type document for the given id
func (s *Service) GetDeviceTypeByID(deviceTypeID string) (*DeviceType, error) {
	path := fmt.Sprintf("%s/%s", _deviceTypesPath, deviceTypeID)

	// Make request
	res, err := s.makeRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("db/GetDeviceTypeByID couch request: %w", err)
	}
	defer res.Body.Close()

	// Check for document not found
	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	// Check for other error codes
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("db/GetDeviceTypeByID got non 200 from couch: Code: %d", res.StatusCode)
	}

	// Read the body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("db/GetDeviceTypeByID read response: %w", err)
	}

	// Unmarshal
	d := DeviceType{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, fmt.Errorf("db/GetDeviceTypeByID json unmarshal: %w", err)
	}

	return &d, nil
}
