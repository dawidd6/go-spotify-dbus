package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/godbus/dbus"
	"github.com/dawidd6/go-spotify-dbus"
)

func main() {
	flag.Parse()
	action := flag.Arg(0)

	conn := getConn()
	var err error
	switch action {
	case "next":
		err = spotify.SendNext(conn)
		return
	case "prev":
		err = spotify.SendPrevious(conn)
		return
	case "play":
		err = spotify.SendPlay(conn)
		return
	case "pause":
		err = spotify.SendPause(conn)
		return
	}
	if err != nil {
		log.Fatalf("%s failed with err: %s", action, err.Error())
	}
	var ps spotify.PlaybackStatus
	var meta *spotify.Metadata
	ps, err = spotify.GetPlaybackStatus(conn)
	if err != nil {
		log.Fatalf("failed getting play/pause status, err: %s", err.Error())
	}
	meta, err = spotify.GetMetadata(conn)
	if err != nil {
		log.Fatalf("failed getting metadata, err: %s", err.Error())
	}

	switch action {
	case "":
		if ps == "Paused" {
			fmt.Printf("▮▮ %s - %s\n", meta.Artist, meta.Title)
		} else {
			fmt.Printf("▶ %s - %s\n", meta.Artist, meta.Title)
		}
	case "info":
		var ld time.Duration
		l := int64(meta.Length)
		ld = time.Duration(time.Duration(l) * time.Microsecond )
		if ps == "Paused" {
			fmt.Printf("▮▮ %s - %s\n\tAlbum: %s\n\tRating: %f\n\tLength: %v\n",
				meta.Artist, meta.Title, meta.Album, meta.AutoRating, ld)
		} else {
			fmt.Printf("▶ %s - %s\n\tAlbum: %s\n\tRating: %f\n\tLength: %v\n",
			meta.Artist, meta.Title, meta.Album, meta.AutoRating, ld)
		}
	}

}

func getConn() *dbus.Conn {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
