package spotify

import (
	"reflect"

	"github.com/godbus/dbus"
)

// Metadata is music player (specifically Spotify) metadata
type Metadata struct {
	Artist      []string `spotify:"xesam:artist"`
	Title       string   `spotify:"xesam:title"`
	Album       string   `spotify:"xesam:album"`
	AlbumArtist []string `spotify:"xesam:albumArtist"`
	AutoRating  float64  `spotify:"xesam:autoRating"`
	DiskNumber  int32    `spotify:"xesam:discNumber"`
	TrackNumber int32    `spotify:"xesam:trackNumber"`
	URL         string   `spotify:"xesam:url"`
	TrackID     string   `spotify:"mpris:trackid"`
	Length      uint64   `spotify:"mpris:length"`
}

// ParseMetadata returns a parsed Metadata struct
func ParseMetadata(variant dbus.Variant) *Metadata {
	metadataMap := variant.Value().(map[string]dbus.Variant)
	metadataStruct := new(Metadata)

	valueOf := reflect.ValueOf(metadataStruct).Elem()
	typeOf := reflect.TypeOf(metadataStruct).Elem()

	for key, val := range metadataMap {
		for i := 0; i < typeOf.NumField(); i++ {
			field := typeOf.Field(i)
			if field.Tag.Get("spotify") == key {
				field := valueOf.Field(i)
				field.Set(reflect.ValueOf(val.Value()))
			}
		}
	}

	return metadataStruct
}
