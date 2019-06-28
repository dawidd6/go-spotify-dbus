package spotify

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus"
)

const (
	sender           = "org.mpris.MediaPlayer2.spotify"
	path             = "/org/mpris/MediaPlayer2"
	member           = "org.mpris.MediaPlayer2.Player"
	playMessage      = member + ".Play"
	pauseMessage     = member + ".Pause"
	playPauseMessage = member + ".PlayPause"
	previousMessage  = member + ".Previous"
	nextMessage      = member + ".Next"
)

// Listeners is a struct of the events we are listening for
type Listeners struct {
	OnMetadata       func(*Metadata)
	OnPlaybackStatus func(PlaybackStatus)
	OnServiceStart   func()
	OnServiceStop    func()
}

// Listen will listen for any changes in PlayPause or metadata from the Spotify app
func Listen(conn *dbus.Conn, listeners *Listeners) {
	currentMetadata := new(Metadata)
	currentPlaybackStatus := StatusUnknown
	newMetadata := new(Metadata)
	newPlaybackStatus := StatusUnknown
	received := new(dbus.Signal)
	signalNameOwnerChanged := "NameOwnerChanged"
	signalPropertiesChanged := "PropertiesChanged"

	args := fmt.Sprintf("sender=%s, path=%s, type=signal, member=PropertiesChanged", sender, path)
	obj := conn.BusObject()
	obj.Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		args,
	)
	args = fmt.Sprintf("type=signal, interface=org.freedesktop.DBus, member=NameOwnerChanged, path=/org/freedesktop/DBus, sender=org.freedesktop.DBus, arg0=%s", sender)
	obj = conn.BusObject()
	obj.Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		args,
	)

	started, err := IsServiceStarted(conn)
	if err != nil {
		return
	} else if started {
		metadata, err := GetMetadata(conn)
		if err != nil {
			return
		}
		currentMetadata = metadata
		listeners.OnMetadata(metadata)

		status, err := GetPlaybackStatus(conn)
		if err != nil {
			return
		}
		currentPlaybackStatus = status
		listeners.OnPlaybackStatus(status)

		listeners.OnServiceStart()
	} else {
		listeners.OnServiceStop()
	}

	channel := make(chan *dbus.Signal, 10)
	conn.Signal(channel)

	for {
		received = <-channel
		name := strings.Split(received.Name, ".")

		switch name[len(name)-1] {
		case signalNameOwnerChanged:
			started, err := IsServiceStarted(conn)
			if err != nil {
				return
			} else if started {
				listeners.OnServiceStart()
				metadata, err := GetMetadata(conn)
				if err != nil {
					return
				}
				currentMetadata = metadata
				listeners.OnMetadata(metadata)
			} else {
				listeners.OnServiceStop()
			}
		case signalPropertiesChanged:
			metadata := received.Body[1].(map[string]dbus.Variant)["Metadata"]
			status := received.Body[1].(map[string]dbus.Variant)["PlaybackStatus"]
			newMetadata = ParseMetadata(metadata)
			newPlaybackStatus = ParsePlaybackStatus(status)

			if currentMetadata.TrackID != newMetadata.TrackID {
				currentMetadata = newMetadata
				listeners.OnMetadata(newMetadata)
			}

			if currentPlaybackStatus != newPlaybackStatus {
				currentPlaybackStatus = newPlaybackStatus
				listeners.OnPlaybackStatus(newPlaybackStatus)
			}
		}
	}
}

// GetMetadata returns the current metadata from the Spotify app
func GetMetadata(conn *dbus.Conn) (*Metadata, error) {
	obj := conn.Object(sender, path)
	property, err := obj.GetProperty(member + ".Metadata")
	if err != nil {
		return nil, err
	}

	return ParseMetadata(property), nil
}

// GetPlaybackStatus returns the current PlayPause status of the Spotify app
func GetPlaybackStatus(conn *dbus.Conn) (PlaybackStatus, error) {
	obj := conn.Object(sender, path)
	property, err := obj.GetProperty(member + ".PlaybackStatus")
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
