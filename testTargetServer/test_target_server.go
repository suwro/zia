package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	cert := flag.String("cert", "", "Certificate file")
	key := flag.String("key", "", "Key file")

	flag.Parse()

	srvHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Println("received request Ip Addr:", req.RemoteAddr)
		_, _ = fmt.Fprintf(rw, "Server Time: %s", time.Now())
	})

	if len(*cert) == 0 || len(*key) == 0 {
		log.Println("Starting test http server at :8081")
		log.Fatal(http.ListenAndServe(":8081", srvHandler))
	} else {
		log.Println("Starting test https server at :8081")
		log.Fatal(http.ListenAndServeTLS(":8081", *cert, *key, srvHandler))
	}
}
