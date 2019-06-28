package main

import (
	"flag"
	"log"

	"github.com/godbus/dbus"
	"github.com/leosunmo/go-spotify-dbus"
)

func main() {
	flag.Parse()
	action := flag.Arg(0)

	conn := getConn()
	var err error
	switch action {
	case "next":
		err = spotify.SendNext(conn)
	case "prev":
		err = spotify.SendPrevious(conn)
	case "play":
		err = spotify.SendPlay(conn)
	case "pause":
		err = spotify.SendPause(conn)
	case "":
		err = spotify.SendPlayPause(conn)
	}
	if err != nil {
		log.Fatalf("%s failed with err: %s", action, err.Error())
	}
}

func getConn() *dbus.Conn {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
