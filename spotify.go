package spotify

import (
	"fmt"
	"github.com/godbus/dbus"
	"strings"
)

const (
	sender = "org.mpris.MediaPlayer2.spotify"
	path   = "/org/mpris/MediaPlayer2"
	member = "org.mpris.MediaPlayer2.Player"

	playMessage           = member + ".Play"
	pauseMessage          = member + ".Pause"
	playPauseMessage      = member + ".PlayPause"
	previousMessage       = member + ".Previous"
	nextMessage           = member + ".Next"
	metadataMessage       = member + ".Metadata"
	playbackStatusMessage = member + ".PlaybackStatus"

	signalNameOwnerChanged  = "NameOwnerChanged"
	signalPropertiesChanged = "PropertiesChanged"

	metadata       = "Metadata"
	playbackStatus = "PlaybackStatus"
)

// Listeners is a struct of the events we are listening for
type Listeners struct {
	OnMetadata       func(*Metadata)
	OnPlaybackStatus func(PlaybackStatus)
	OnServiceStart   func()
	OnServiceStop    func()
}

// Listen will listen for any changes in PlayPause or metadata from the Spotify app
//
// This function is blocking
func Listen(conn *dbus.Conn, errors chan error, listeners *Listeners) {
	var (
		currentMetadata       = new(Metadata)
		currentPlaybackStatus = StatusUnknown
		channel               = make(chan *dbus.Signal, 10)
	)

	// Register channel receiving signals
	conn.Signal(channel)

	// Watch for signals about metadata changes
	args := fmt.Sprintf("sender=%s, path=%s, type=signal, member=PropertiesChanged", sender, path)
	conn.BusObject().Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		args,
	)

	// Watch for signals about service status changes
	args = fmt.Sprintf("type=signal, interface=org.freedesktop.DBus, member=NameOwnerChanged, path=/org/freedesktop/DBus, sender=org.freedesktop.DBus, arg0=%s", sender)
	conn.BusObject().Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		args,
	)

	// Initially check if service is up
	started, err := IsServiceStarted(conn)
	if err != nil {
		errors <- err
	}

	// If up, then initially get metadata and playback status
	if started {
		newMetadata, err := GetMetadata(conn)
		if err != nil {
			errors <- err
		}

		newPlaybackStatus, err := GetPlaybackStatus(conn)
		if err != nil {
			errors <- err
		}

		currentMetadata = newMetadata
		currentPlaybackStatus = newPlaybackStatus

		listeners.OnServiceStart()
		listeners.OnPlaybackStatus(newPlaybackStatus)
		listeners.OnMetadata(newMetadata)
	} else {
		listeners.OnServiceStop()
	}

	// Listen for changes for ever
	for {
		// Wait for signal
		signal := <-channel

		// New metadata and playback status received
		if strings.HasSuffix(signal.Name, signalPropertiesChanged) {
			newMetadata := ParseMetadata(signal.Body[1].(map[string]dbus.Variant)[metadata])
			newPlaybackStatus := ParsePlaybackStatus(signal.Body[1].(map[string]dbus.Variant)[playbackStatus])

			if currentMetadata.TrackID != newMetadata.TrackID {
				currentMetadata = newMetadata
				listeners.OnMetadata(newMetadata)
			}

			if currentPlaybackStatus != newPlaybackStatus {
				currentPlaybackStatus = newPlaybackStatus
				listeners.OnPlaybackStatus(newPlaybackStatus)
			}
		}

		// Service status changed (Spotify was closed or opened, probably)
		if strings.HasSuffix(signal.Name, signalNameOwnerChanged) {
			started, err := IsServiceStarted(conn)
			if err != nil {
				errors <- err
				return
			}

			if started {
				listeners.OnServiceStart()
			} else {
				listeners.OnServiceStop()
			}
		}
	}
}

// GetMetadata returns the current metadata from the Spotify app
func GetMetadata(conn *dbus.Conn) (*Metadata, error) {
	obj := conn.Object(sender, path)
	property, err := obj.GetProperty(metadataMessage)
	if err != nil {
		return nil, err
	}

	return ParseMetadata(property), nil
}

// GetPlaybackStatus returns the current Play/Pause status of the Spotify app
// Status will be "Playing", "Paused" or "Unknown"
func GetPlaybackStatus(conn *dbus.Conn) (PlaybackStatus, error) {
	obj := conn.Object(sender, path)
	property, err := obj.GetProperty(playbackStatusMessage)
	if err != nil {
		return StatusUnknown, err
	}

	return ParsePlaybackStatus(property), nil
}

// IsServiceStarted checks if the Spotify app is running
func IsServiceStarted(conn *dbus.Conn) (bool, error) {
	started := false

	err := conn.Object(
		"org.freedesktop.DBus",
		"/org/freedesktop/DBus",
	).Call(
		"org.freedesktop.DBus.NameHasOwner",
		0,
		sender,
	).Store(
		&started,
	)
	if err != nil {
		return false, err
	}

	return started, nil
}

// SendPlay sends a "Play" message to the Spotify app.
// Returns error if anything goes wrong.
// If the Spotify app is not running, return nil
func SendPlay(conn *dbus.Conn) error {
	started, err := IsServiceStarted(conn)
	if err != nil {
		return err
	} else if started {
		obj := conn.Object(sender, path)
		c := obj.Call(playMessage, 0)
		if c.Err != nil {
			return c.Err
		}
	}
	return nil
}

// SendPause sends a "Pause" message to the Spotify app.
// Returns error if anything goes wrong.
// If the Spotify app is not running, return nil
func SendPause(conn *dbus.Conn) error {
	started, err := IsServiceStarted(conn)
	if err != nil {
		return err
	} else if started {
		obj := conn.Object(sender, path)
		c := obj.Call(pauseMessage, 0)
		if c.Err != nil {
			return c.Err
		}
	}
	return nil
}

// SendPlayPause sends a "PlayPause" message to the Spotify app.
// Returns error if anything goes wrong.
// If the Spotify app is not running, return nil
func SendPlayPause(conn *dbus.Conn) error {
	started, err := IsServiceStarted(conn)
	if err != nil {
		return err
	} else if started {
		obj := conn.Object(sender, path)
		c := obj.Call(playPauseMessage, 0)
		if c.Err != nil {
			return c.Err
		}
	}
	return nil
}

// SendNext sends a "Next" message to the Spotify app.
// Returns error if anything goes wrong.
// If the Spotify app is not running, return nil
func SendNext(conn *dbus.Conn) error {
	started, err := IsServiceStarted(conn)
	if err != nil {
		return err
	} else if started {
		obj := conn.Object(sender, path)
		c := obj.Call(nextMessage, 0)
		if c.Err != nil {
			return c.Err
		}
	}
	return nil
}

// SendPrevious sends a "Previous" message to the Spotify app.
// Returns error if anything goes wrong.
// If the Spotify app is not running, return nil
func SendPrevious(conn *dbus.Conn) error {
	started, err := IsServiceStarted(conn)
	if err != nil {
		return err
	} else if started {
		obj := conn.Object(sender, path)
		c := obj.Call(previousMessage, 0)
		if c.Err != nil {
			return c.Err
		}
	}
	return nil
}
