package helpers

import (
	"github.com/byuoitav/av-api/base"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/uapi-translator/structs"
)

// UAPItoAV takes the Resource structure of the UAPI and translates it into the AV-API PublicRoom structure.
func UAPItoAV(resource structs.Resource) (pubRoom base.PublicRoom, ne *nerr.E) {
	return pubRoom, ne
}

// AVtoUAPI takes the AV-API PublicRoom structure and translates it into the UAPI Resource structure.
func AVtoUAPI(pubRoom base.PublicRoom) (resource structs.Resource, ne *nerr.E) {
	return resource, ne
}
