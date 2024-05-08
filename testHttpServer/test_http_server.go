package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	srvHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Println("received request Ip Addr:", req.RemoteAddr)
		_, _ = fmt.Fprintf(rw, "Server Time: %s", time.Now())
	})

	log.Println("Starting server at :8081...")
	log.Fatal(http.ListenAndServe(":8081", srvHandler))
}
