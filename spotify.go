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

var stop = make(chan int)

func Stop() {
	stop <- 1
}

func WaitForPropertiesChanges(conn *dbus.Conn, onMetadata func(*Metadata), onStatus func(PlaybackStatus)) {
	currentMetadata := new(Metadata)
	currentPlaybackStatus := StatusUnknown
	newMetadata := new(Metadata)
	newPlaybackStatus := StatusUnknown
	received := new(dbus.Signal)

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

	channel := make(chan *dbus.Signal, 10)
	conn.Signal(channel)

	for {
		select {
		case received = <-channel:
		case <-stop:
			return
		}

		if strings.Contains(received.Name, "NameOwnerChanged") {
			newMetadata = new(Metadata)
			newPlaybackStatus = StatusUnknown
		} else {
			metadata := received.Body[1].(map[string]dbus.Variant)["Metadata"]
			status := received.Body[1].(map[string]dbus.Variant)["PlaybackStatus"]
			newMetadata = ParseMetadata(metadata)
			newPlaybackStatus = ParsePlaybackStatus(status)
		}

		if currentMetadata.TrackID != newMetadata.TrackID {
			currentMetadata = newMetadata
			onMetadata(newMetadata)
		}

		if currentPlaybackStatus != newPlaybackStatus {
			currentPlaybackStatus = newPlaybackStatus
			onStatus(newPlaybackStatus)
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
