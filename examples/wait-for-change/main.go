package main

import (
	"encoding/json"
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
			// Marshal struct to indented json for better readability
			jason, err := json.MarshalIndent(metadata, "", "  ")
			if err != nil {
				errors <- err
			}

			log.Println("metadata:", string(jason))
		},
		OnPlaybackStatus: func(status spotify.PlaybackStatus) {
			log.Println("status:", status)
		},
		OnServiceStart: func() {
			log.Println("service: Start")
		},
		OnServiceStop: func() {
			log.Println("service: Stop")
		},
	}

	// Print errors in background
	go func() {
		for {
			err := <-errors
			log.Println("error:", err)
		}
	}()

	// Listen for changes, blocking further execution
	spotify.Listen(conn, errors, listeners)
}
