package main

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var (
	addr = ":30001"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/hello", hello)
	log.Infoln("listen addr ", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalln(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "HOST = "+r.Header.Get("HOST"))
}

func index(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "index")
}
