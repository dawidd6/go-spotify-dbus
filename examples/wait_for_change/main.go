package main

import (
	"fmt"
	"github.com/dawidd6/go-spotify-dbus"
	"log"
)

func main() {
	spot, err := spotify.New()
	if err != nil {
		log.Fatal(err)
	}

	spot.OnMetadata = func(metadata *spotify.Metadata) {
		fmt.Printf("METADATA: %+v\n", metadata)
	}
	spot.OnStatus = func(status spotify.PlaybackStatus) {
		fmt.Printf("STATUS: %s\n", status)
	}

	spot.WaitForPropertiesChanges()
}
