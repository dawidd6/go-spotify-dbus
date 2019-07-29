package spotify

import (
	"github.com/godbus/dbus"
	"strings"
)

// PlaybackStatus is a PlayPause status of a music player
type PlaybackStatus string

const (
	// StatusPlaying is the playing state
	StatusPlaying PlaybackStatus = "Playing"
	// StatusPaused is the paused state
	StatusPaused PlaybackStatus = "Paused"
	// StatusUnknown is an unknown music player state
	StatusUnknown PlaybackStatus = "Unknown"
)

// parsePlaybackStatus parses the current PlayPause status
func parsePlaybackStatus(variant dbus.Variant) PlaybackStatus {
	status := strings.Trim(variant.String(), "\"")

	switch status {
	case string(StatusPlaying):
		return StatusPlaying
	case string(StatusPaused):
		return StatusPaused
	default:
		return StatusUnknown
	}
}
