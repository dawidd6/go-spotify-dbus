package spotify

import (
	"github.com/godbus/dbus"
	"strings"
)

const (
	StatusPlaying PlaybackStatus = "Playing"
	StatusPaused  PlaybackStatus = "Paused"
	StatusUnknown PlaybackStatus = "Unknown"
)

type PlaybackStatus string

func ParsePlaybackStatus(variant dbus.Variant) PlaybackStatus {
	if strings.Contains(variant.String(), "Playing") {
		return StatusPlaying
	}
	if strings.Contains(variant.String(), "Paused") {
		return StatusPaused
	}

	return StatusUnknown
}
