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

	listeners := &spotify.Listeners{
		OnMetadata: func(metadata *spotify.Metadata) {
			fmt.Println("metadata: ", metadata)
		},
		OnPlaybackStatus: func(status spotify.PlaybackStatus) {
			fmt.Println("status: ", status)
		},
		OnServiceStart: func() {
			fmt.Println("start")
		},
		OnServiceStop: func() {
			fmt.Println("stop")
		},
	}

	spotify.Listen(conn, listeners)
}
