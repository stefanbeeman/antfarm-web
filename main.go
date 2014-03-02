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
	world        af.World
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
	// world
	world = af.MakeWorld("../antfarm-data", 20, 20, 0)
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
	ns.Emit("world", world)
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

	// Handle routes
	// http.Handle("/", http.FileServer(staticRoutes))

	// Listen on hostname:port
	// fmt.Printf("Listening on %s:%d...\n", hostname, port)
	// err := http.ListenAndServe(fmt.Sprintf("%s:%d", hostname, port), nil)
	// if err != nil {
	//	log.Fatal("Error: ", err)
	// }

	sock_config := &socketio.Config{}
	sock_config.HeartbeatTimeout = 2
	sock_config.ClosingTimeout = 4

	sio := socketio.NewSocketIOServer(sock_config)

	// Handler for new connections, also adds socket.io event handlers
	sio.On("connect", onConnect)
	sio.On("disconnect", onDisconnect)
	ticker := time.NewTicker(time.Millisecond * 1000)
	go func() {
		for t := range ticker.C {
			world.RunFor(1)
			sio.Broadcast("world", world, t)
		}
	}()

	//this will serve a http static file server
	sio.Handle("/", http.FileServer(staticRoutes))
	//startup the server
	http.ListenAndServe(":9000", sio)
}
