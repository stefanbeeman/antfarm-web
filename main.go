package main

import (
	"flag"
	"fmt"
	"github.com/googollee/go-socket.io"
	"github.com/stefanbeeman/antfarm"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	hostname     string
	port         int
	topStaticDir string
	game         antfarm.Game
)

func init() {
	// Flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [default_static_dir]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&hostname, "h", "localhost", "hostname")
	flag.IntVar(&port, "p", 8080, "port")
	flag.StringVar(&topStaticDir, "static_dir", "", "static directory in addition to default static directory")
	game = antfarm.MakeGame("../antfarm-data", 20, 20, 1)
}

func appendStaticRoute(sr StaticRoutes, dir string) StaticRoutes {
	if _, err := os.Stat(dir); err != nil {
		log.Fatal(err)
	}
	return append(sr, http.Dir(dir))
}

type StaticRoutes []http.FileSystem

func (sr StaticRoutes) Open(name string) (f http.File, err error) {
	for _, s := range sr {
		if f, err = s.Open(name); err == nil {
			f = disabledDirListing{f}
			return
		}
	}
	return
}

type disabledDirListing struct {
	http.File
}

func (f disabledDirListing) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func onConnect(ns *socketio.NameSpace) {
	fmt.Println("connected:", ns.Id())
	ns.Emit("game", game.Display())
}

func onDisconnect(ns *socketio.NameSpace) {
	fmt.Println("disconnected:", ns.Id())
}

func main() {
	// Parse flags
	flag.Parse()
	staticDir := flag.Arg(0)

	// Setup static routes
	staticRoutes := make(StaticRoutes, 0)
	if topStaticDir != "" {
		staticRoutes = appendStaticRoute(staticRoutes, topStaticDir)
	}
	if staticDir == "" {
		staticDir = "./"
	}
	staticRoutes = appendStaticRoute(staticRoutes, staticDir)

	sock_config := &socketio.Config{}
	sock_config.HeartbeatTimeout = 2
	sock_config.ClosingTimeout = 4

	sio := socketio.NewSocketIOServer(sock_config)

	sio.On("connect", onConnect)
	sio.On("disconnect", onDisconnect)

	ticker := time.NewTicker(time.Millisecond * 100)
	go func() {
		for t := range ticker.C {
			game.RunFor(1)
			sio.Broadcast("game", game.Display(), t)
		}
	}()

	//this will serve a http static file server
	sio.Handle("/", http.FileServer(staticRoutes))
	//startup the server
	http.ListenAndServe(":9000", sio)
}
