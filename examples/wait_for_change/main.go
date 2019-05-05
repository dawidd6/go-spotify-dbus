package main

import (
	"fmt"
	"github.com/dawidd6/go-spotify-dbus"
	"github.com/godbus/dbus"
	"log"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}

	spotify.WaitForPropertiesChanges(
		conn,
		func(metadata *spotify.Metadata) {
			fmt.Println(metadata)
		},
		func(status spotify.PlaybackStatus) {
			fmt.Println(status)
		},
	)
}
