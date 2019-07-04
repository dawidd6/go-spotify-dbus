package spotify

import (
	"strings"

	"github.com/godbus/dbus"
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

// ParsePlaybackStatus parses the current PlayPause status
func ParsePlaybackStatus(variant dbus.Variant) PlaybackStatus {
	if strings.Contains(variant.String(), "Playing") {
		return StatusPlaying
	}
	if strings.Contains(variant.String(), "Paused") {
		return StatusPaused
	}

	return StatusUnknown
}
