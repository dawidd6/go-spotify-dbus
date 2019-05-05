package spotify

import (
	"errors"
	"fmt"
	"github.com/godbus/dbus"
)

const (
	sender = "org.mpris.MediaPlayer2.spotify"
	dest   = sender
	path   = "/org/mpris/MediaPlayer2"
	member = "org.mpris.MediaPlayer2.Player"
)

type Spotify struct {
	conn                  *dbus.Conn
	OnMetadata            func(metadata *Metadata)
	OnStatus              func(status PlaybackStatus)
	CurrentMetadata       *Metadata
	CurrentPlaybackStatus PlaybackStatus
}

func New() (*Spotify, error) {
	err := errors.New("")
	spotify := new(Spotify)

	spotify.conn, err = dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	spotify.CurrentMetadata, err = spotify.GetMetadata()
	if err != nil {
		return nil, err
	}

	spotify.CurrentPlaybackStatus, err = spotify.GetPlaybackStatus()
	if err != nil {
		return nil, err
	}

	return spotify, nil
}

func (spotify *Spotify) WaitForPropertiesChanges() {
	args := fmt.Sprintf("sender=%s, path=%s, type=signal, member=PropertiesChanged", sender, path)
	obj := spotify.conn.BusObject()
	obj.Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		args,
	)

	channel := make(chan *dbus.Signal)
	spotify.conn.Signal(channel)

	for {
		received := <-channel
		metadata := received.Body[1].(map[string]dbus.Variant)["Metadata"]
		status := received.Body[1].(map[string]dbus.Variant)["PlaybackStatus"]

		newMetadata := ParseMetadata(metadata)
		newStatus := ParsePlaybackStatus(status)

		if spotify.CurrentMetadata.TrackID != newMetadata.TrackID {
			spotify.CurrentMetadata = newMetadata
			spotify.OnMetadata(newMetadata)
		}

		if spotify.CurrentPlaybackStatus != newStatus {
			spotify.CurrentPlaybackStatus = newStatus
			spotify.OnStatus(newStatus)
		}
	}
}

func (spotify *Spotify) GetMetadata() (*Metadata, error) {
	obj := spotify.conn.Object(dest, path)
	property, err := obj.GetProperty(member + ".Metadata")
	if err != nil {
		return nil, err
	}

	return ParseMetadata(property), nil
}

func (spotify *Spotify) GetPlaybackStatus() (PlaybackStatus, error) {
	obj := spotify.conn.Object(dest, path)
	property, err := obj.GetProperty(member + ".PlaybackStatus")
	if err != nil {
		return StatusUnknown, err
	}

	return ParsePlaybackStatus(property), nil
}
