package structs

//Rooms
type Room struct {
	roomId   string `json:"av_room_id"`
	roomNum  string `json:"room_number"`
	bldgAbbr string `json:"building_abbreviation"`
}

type RoomDevices struct {
	displays []Display     `json:"av_displays"`
	outputs  []AudioOutput `json:"av_audio_outputs"`
	inputs   []Input       `json:"av_inputs"`
}

//Devices
type Device struct {
	deviceId   string `json:"av_device_id"`
	deviceName string `json:"av_device_name"`
	deviceType string `json:"av_device_type"`
	bldgAbbr   string `json:"building_abbreviation"`
	roomNum    string `json:"room_number"`
}

type DeviceProperty struct {
	name  string `json:"av_device_property_name"`
	value string `json:"av_device_property_value"`
}

type DeviceStateAttribute struct {
	name  string `json:"av_device_state_attribute_name"`
	value string `json:"av_device_state_attribute_value"`
}

//Inputs
type Input struct {
	deviceId   string   `json:"av_device_id"`
	roomNum    string   `json:"room_number"`
	bldgAbbr   string   `json:"building_abbreviation"`
	deviceType string   `json:"av_device_type"`
	outputs    []string `json:"av_outputs"`
}

//Displays
type Display struct {
	displayId string `json:"av_display_id"`
	roomNum   string `json:"room_number"`
	bldgAbbr  string `json:"building_abbreviation"`
}

type DisplayConfig struct {
	devices []Device `json:"av_devices"`
	inputs  []Input  `json:"av_inputs"`
}

type DisplayState struct {
	powered bool   `json:"av_display_powered"`
	blanked bool   `json:"av_display_blanked"`
	input   string `json:"av_display_input"`
}

//Audio Outputs
type AudioOutput struct {
	outputId   string `json:"av_audio_output_id"`
	roomNum    string `json:"room_number"`
	bldgAbbr   string `json:"building_abbreviation"`
	deviceType string `json:"av_device_type"`
}

type AudioOutputState struct {
	volume int  `json:"av_audio_output_volume_level"`
	muted  bool `json:"av_audio_output_muted"`
}
