package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var version = "0.1.0"

func main() {
	cert := flag.String("cert", "", "Certificate file")
	key := flag.String("key", "", "Key file")
	port := flag.Int("port", 8081, "Tcp port to listen on")
	ip := flag.String("ip", "", "IP address to bind to")
	flag.Parse()

	srvHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Println("received request Ip Addr:", req.RemoteAddr)
		_, _ = fmt.Fprintf(rw, "Server Version: %s Time: %s", version, time.Now())
	})

	srvAddr := fmt.Sprintf("%s:%d", *ip, *port)

	if len(*cert) == 0 || len(*key) == 0 {
		log.Println("Starting http server at:", srvAddr)
		log.Fatal(http.ListenAndServe(srvAddr, srvHandler))
	} else {
		log.Println("Starting https server at:", srvAddr)
		log.Fatal(http.ListenAndServeTLS(srvAddr, *cert, *key, srvHandler))
	}
}
