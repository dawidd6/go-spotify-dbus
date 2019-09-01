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

	// Get default structure of listeners
	listeners := spotify.NewListeners()

	// Define selected listeners that are fired upon an event
	listeners.OnMetadata = func(metadata *spotify.Metadata) {
		// Marshal struct to indented json for better readability
		jason, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			log.Println("error:", err)
			return
		}

		log.Println("metadata:", string(jason))
	}
	listeners.OnPlaybackStatus = func(status spotify.PlaybackStatus) {
		log.Println("status:", status)
	}
	listeners.OnServiceStart = func() {
		log.Println("service: Start")
	}
	listeners.OnServiceStop = func() {
		log.Println("service: Stop")
	}
	listeners.OnError = func(e error) {
		log.Println("error", err)
	}

	// Listen for changes, blocking further execution
	spotify.Listen(conn, listeners)
}
