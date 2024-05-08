package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/suwro/zia/src/proxy"
	"golang.org/x/crypto/acme"
)

func main() {
	cfgFile := flag.String("c", "config/config.json", "config file")

	flag.Parse()
	proxy.InitConfig(*cfgFile)
	proxy.Version = "0.0.4 beta"

	log.Println(proxy.Au.White("Zia"), "reverse proxy ver:", proxy.Au.Yellow(proxy.Version))

	// Static
	fs := http.FileServer(http.Dir("static"))
	h := &proxy.BaseHandle{}

	// Http Handler
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", h)

	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", proxy.Cfg.Service.IP, proxy.Cfg.Service.Port),
		ReadTimeout: 5 * time.Second,
	}

	// https scheme
	if proxy.Cfg.Service.SSL {
		log.Println("Starting", proxy.Au.Red("https"), "proxy:", proxy.Au.Yellow(proxy.Cfg.Service.IP), "port:", proxy.Au.Green(proxy.Cfg.Service.Port))
		/*
		   / AUTO TLS
		   			autoTLSManager := autocert.Manager{
		   				Prompt: autocert.AcceptTOS,
		   				// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
		   				Cache:      autocert.DirCache("config/cert"),
		   				HostPolicy: autocert.HostWhitelist(tlsDomainName),
		   				//Email:      "dacian.stanciu@just.ro",
		   			}
		*/
		server.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			//GetCertificate: autoTLSManager.GetCertificate,
			NextProtos: []string{acme.ALPNProto},
		}

		log.Fatal(server.ListenAndServeTLS("ziacert.pem", "ziaca.key"))
		return
	}
	log.Println("Starting http proxy:", proxy.Au.Yellow(proxy.Cfg.Service.IP), "port:", proxy.Au.Green(proxy.Cfg.Service.Port))
	log.Fatal(server.ListenAndServe())
}
