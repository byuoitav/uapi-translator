package models

//Rooms
type Room struct {
	RoomID   string `json:"av_room_id"`
	RoomNum  string `json:"room_number"`
	BldgAbbr string `json:"building_abbreviation"`
}

type RoomDevices struct {
	Displays []Display     `json:"av_displays"`
	Outputs  []AudioOutput `json:"av_audio_outputs"`
	Inputs   []Input       `json:"av_inputs"`
}

//Devices
type Device struct {
	DeviceID   string `json:"av_device_id"`
	DeviceName string `json:"av_device_name"`
	DeviceType string `json:"av_device_type"`
	BldgAbbr   string `json:"building_abbreviation"`
	RoomNum    string `json:"room_number"`
}

type DeviceProperty struct {
	Name  string `json:"av_device_property_name"`
	Value string `json:"av_device_property_value"`
}

type DeviceStateAttribute struct {
	Name  string `json:"av_device_state_attribute_name"`
	Value string `json:"av_device_state_attribute_value"`
}

//Inputs
type Input struct {
	DeviceID   string   `json:"av_device_id"`
	RoomNum    string   `json:"room_number"`
	BldgAbbr   string   `json:"building_abbreviation"`
	DeviceType string   `json:"av_device_type"`
	Outputs    []string `json:"av_outputs"`
}

//Displays
type Display struct {
	DisplayID string `json:"av_display_id"`
	RoomNum   string `json:"room_number"`
	BldgAbbr  string `json:"building_abbreviation"`
}

type DisplayConfig struct {
	Devices []Device `json:"av_devices"`
	Inputs  []Input  `json:"av_inputs"`
}

type DisplayState struct {
	Powered bool   `json:"av_display_powered"`
	Blanked bool   `json:"av_display_blanked"`
	Input   string `json:"av_display_input"`
}

//Audio Outputs
type AudioOutput struct {
	OutputID   string `json:"av_audio_output_id"`
	RoomNum    string `json:"room_number"`
	BldgAbbr   string `json:"building_abbreviation"`
	DeviceType string `json:"av_device_type"`
}

type AudioOutputState struct {
	Volume int  `json:"av_audio_output_volume_level"`
	Muted  bool `json:"av_audio_output_muted"`
}
