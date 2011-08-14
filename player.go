package main

import (
	"flag"
	"http"
	"json"
	"os"
)

const (
	filePrefix = "/f/"
)

var (
	addr = flag.String("http", ":8080", "http listen address")
	root = flag.String("root", "/home/ton/Music/", "music root")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", Index)
	http.HandleFunc(filePrefix, File)
	http.ListenAndServe(":8080", nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "index.html")
}

func File(w http.ResponseWriter, r *http.Request) {
	fn := *root + r.URL.Path[len(filePrefix):]
	fi, err := os.Stat(fn)
	if err != nil {
		http.Error(w, err.String(), http.StatusNotFound)
		return
	}
	if fi.IsDirectory() {
		serveDirectory(fn, w, r)
		return
	}
	http.ServeFile(w, r, fn)
}

func serveDirectory(fn string, w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err, ok := recover().(os.Error); ok {
			http.Error(w, err.String(), http.StatusInternalServerError)
		}
	}()
	d, err := os.Open(fn)
	if err != nil {
		panic(err)
	}
	files, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}
	j := json.NewEncoder(w)
	if err := j.Encode(files); err != nil {
		panic(err)
	}
}

