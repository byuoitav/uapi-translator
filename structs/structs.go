package structs

import "strings"

// Resource is the base level object returned by the UAPI.
type Resource struct {
	Links    map[string]Link `json:"links,omitempty"`
	Metadata Metadata        `json:"metadata,omitempty"`
	Basic    interface{}     `json:"basic,omitempty"`
	State    interface{}     `json:"av_state,omitempty"`
	Config   interface{}     `json:"av_config,omitempty"`
}

// SubResource is an object that helps to comprise a Resource.
type SubResource struct {
	Links        map[string]Link `json:"links,omitempty"`
	Metadata     Metadata        `json:"metadata,omitempty"`
	Building     interface{}     `json:"building,omitempty"`
	Room         interface{}     `json:"room,omitempty"`
	Displays     interface{}     `json:"displays,omitempty"`
	AudioDevices interface{}     `json:"audio_devices,omitempty"`
}

// Link contains information about accessing the Resource.
type Link struct {
	Rel    string `json:"rel,omitempty"`
	Href   string `json:"href,omitempty"`
	Method string `json:"method,omitempty"`
}

// Metadata contains high level metadata about the Resource or SubResource
type Metadata struct {
	ValidationResponse interface{} `json:"validation_response,omitempty"`
	// ValidationInformation []string           `json:"validation_information,omitempty"`
	// Cache                 Cache              `json:"cache,omitempty"`
	// Restricted            *bool              `json:"restricted,omitempty"`
	FieldSetsReturned  []string `json:"field_sets_returned,omitempty"`
	FieldSetsAvailable []string `json:"field_sets_available,omitempty"`
	FieldSetsDefault   []string `json:"field_sets_default,omitempty"`
}

// ValidationResponse has information about the request.
type ValidationResponse struct {
	Code    *int   `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Cache contains a DateTime about the Resource if it was cached.
type Cache struct {
	// DateTime string `json:"date_time,omitempty"`
}

// Property is an attribute of a Resource or SubResource.
type Property struct {
	Type        string                 `json:"api_type,omitempty"`
	Key         bool                   `json:"key,omitempty"`
	Value       string                 `json:"value,omitempty"`
	Object      interface{}            `json:"object,omitempty"`
	ValueArray  []string               `json:"value_array,omitempty"`
	ObjectArray map[string]interface{} `json:"object_array,omitempty"`
	// Description     string        `json:"description,omitempty"`
	// DisplayLabel    string        `json:"display_label,omitempty"`
	// Domain          string        `json:"domain,omitempty"`
	// LongDescription string        `json:"long_description,omitempty"`
	// RelatedResource string        `json:"related_resource,omitempty"`
}

// ReachableRoomConfig represents the amount of data we will report back about the room.
type ReachableRoomConfig struct {
	Room
	InputReachability map[string][]string `json:"input_reachability"`
}

// Room is the filtered Room struct for returning from the UAPI.
type Room struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Designation string   `json:"designation"`
	Devices     []Device `json:"devices,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Device is the filtered Device struct for returning from the UAPI.
type Device struct {
	ID          string     `json:"_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	DisplayName string     `json:"display_name"`
	Type        DeviceType `json:"type,omitempty"`
	Roles       []Role     `json:"roles"`
	Tags        []string   `json:"tags,omitempty"`
}

// DeviceType is the filtered DeviceType struct for returning from the UAPI.
type DeviceType struct {
	ID          string `json:"_id"`
	Description string `json:"description,omitempty"`
}

// Role is the filtered Role struct for returning from the UAPI.
type Role struct {
	ID string `json:"_id"`
}

// HasRole checks to see if the given device has the given role.
func (d *Device) HasRole(role string) bool {
	role = strings.ToLower(role)
	for i := range d.Roles {
		if strings.EqualFold(strings.ToLower(d.Roles[i].ID), role) {
			return true
		}
	}
	return false
}
