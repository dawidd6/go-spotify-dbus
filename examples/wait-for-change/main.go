package main

import (
	"log"

	"github.com/dawidd6/go-spotify-dbus"
	"github.com/godbus/dbus"
)

func main() {
	// Establish connection with session dbus
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}

	// Make channel with errors
	errors := make(chan error, 5)

	// Define a bunch of listeners that are fired upon an event
	listeners := &spotify.Listeners{
		OnMetadata: func(metadata *spotify.Metadata) {
			log.Println("metadata:", metadata)
		},
		OnPlaybackStatus: func(status spotify.PlaybackStatus) {
			log.Println("status:", status)
		},
		OnServiceStart: func() {
			log.Println("service: start")
		},
		OnServiceStop: func() {
			log.Println("service: stop")
		},
	}

	// Print errors in background
	go func() {
		err := <-errors
		log.Println("error:", err)
	}()

	// Listen for changes, blocking further execution
	spotify.Listen(conn, errors, listeners)
}
