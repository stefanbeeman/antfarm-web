package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/stefanbeeman/antfarm"
	"log"
	"net/http"
	"os"
)

var (
	hostname     string
	port         int
	topStaticDir string
	myWorld      antfarm.World
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
	myWorld = antfarm.MakeWorld("../antfarm-data", 12, 12)
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

type Profile struct {
	Name    string
	Hobbies []string
}

func getWorld(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(myWorld)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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
	http.Handle("/", http.FileServer(staticRoutes))
	http.HandleFunc("/world", getWorld)

	// Listen on hostname:port
	fmt.Printf("Listening on %s:%d...\n", hostname, port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", hostname, port), nil)
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
