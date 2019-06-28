package main

import (
	"fmt"
	"log"

	"github.com/godbus/dbus"
	"github.com/leosunmo/go-spotify-dbus"
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
