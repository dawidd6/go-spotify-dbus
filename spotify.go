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
)

type Listeners struct {
	OnMetadata       func(*Metadata)
	OnPlaybackStatus func(PlaybackStatus)
	OnServiceStart   func()
	OnServiceStop    func()
}

func Listen(conn *dbus.Conn, listeners *Listeners) {
	currentMetadata := new(Metadata)
	currentPlaybackStatus := StatusUnknown
	newMetadata := new(Metadata)
	newPlaybackStatus := StatusUnknown
	received := new(dbus.Signal)
	signalNameOwnerChanged := "NameOwnerChanged"
	signalPropertiesChanged := "PropertiesChanged"
	started := false

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
		return
	}

	channel := make(chan *dbus.Signal, 10)
	conn.Signal(channel)

	for {
		received = <-channel
		name := strings.Split(received.Name, ".")

		switch name[len(name)-1] {
		case signalNameOwnerChanged:
			started = !started
			if started {
				listeners.OnServiceStart()
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

/*
func GetMetadata(conn *dbus.Conn) (*Metadata, error) {
	obj := conn.Object(dest, path)
	property, err := obj.GetProperty(member + ".Metadata")
	if err != nil {
		return nil, err
	}

	return ParseMetadata(property), nil
}

func GetPlaybackStatus(conn *dbus.Conn) (PlaybackStatus, error) {
	obj := conn.Object(dest, path)
	property, err := obj.GetProperty(member + ".PlaybackStatus")
	if err != nil {
		return StatusUnknown, err
	}

	return ParsePlaybackStatus(property), nil
}
*/
