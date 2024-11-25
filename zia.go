package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/acme"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

var versiune = "0.3.3"

func main() {
	cert := flag.String("cert", "", "Certificate file")
	key := flag.String("key", "", "Key file")
	domain := flag.String("domain", "", "Domain name to use for the certificate")
	port := flag.Int("port", 8080, "Port to listen on")
	ssl := flag.Bool("ssl", true, "Use SSL/TLS default true")
	ver := flag.Bool("version", false, "Show version and exit")
	targetList := flag.String("targets", "", "List of targets for proxy, comma separated")
	timeout := flag.Int("timeout", 0, "Timeout for proxy in seconds, 0 no timeout")
	stdout := flag.Bool("stdout", false, "Use stdout instead /var/log/zia/<domain>/acces_<port>.log")

	flag.Parse()
	log.Println("Versiune:", versiune, "SSL:", *ssl)
	if *ver {
		os.Exit(1)
	}

	// verifica domeniu
	if len(*domain) == 0 {
		log.Fatal("Domain name is required")
	}

	if len(*cert) > 0 && len(*key) > 0 {
		*ssl = true
		log.Println("Certificate and key files provided")
	}

	// verifica lista de tinte proxy
	lista, err := parseTargets(*targetList)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Domain:", *domain, "Port:", *port, "Targets:", lista)

	// Logs
	var logConfig = middleware.LoggerConfig{
		Format: "${time_rfc3339}\t${remote_ip}\t${method}\t${uri}\t${status} ${error}\n",
	}

	if !*stdout {
		// log in fisier in loc de standard
		logFileName := fmt.Sprintf("/var/log/zia/%s/acces_%d.log", *domain, *port)

		// Salveaza log-urile in logs.txt
		f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error generating log file: %v", err)
		}

		// activeaza magia log-urilor
		defer f.Close()
		log.SetOutput(f)
		logConfig.Output = f
	}

	// Server http proxy
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(logConfig))

	// TLS Transport proxy
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			//MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		},
	}

	// Proxy
	targets, err := addTarget(lista)
	if err != nil {
		log.Fatal(err.Error())
	}

	balancer := middleware.NewRoundRobinBalancer(targets)
	e.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer:  balancer,
		Transport: transport,
	}))

	// http settings
	s := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: e, // Echo instance handler
	}

	// timeout
	if (*timeout > 0) && (*timeout < 3600) {
		s.ReadTimeout = time.Duration(*timeout) * time.Second
	} else if *timeout > 3600 {
		log.Fatal("Timeout too high, set it to 0 for no timeout")
	}

	// server https
	if *ssl {
		log.Println("SSL/TLS enabled")

		s.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}

		// Given cert and key
		if len(*cert) > 0 || len(*key) > 0 {
			log.Fatal(s.ListenAndServeTLS(*cert, *key))
		} else {
			// let's encrypt certificate
			certPath := filepath.Join("config", "cert")
			err = os.MkdirAll(certPath, os.ModePerm)
			if err != nil {
				log.Fatal(err.Error())
			}

			// tls settings
			autoTLSManager := autocert.Manager{
				Prompt: autocert.AcceptTOS,
				// Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
				Cache:      autocert.DirCache(certPath),
				HostPolicy: autocert.HostWhitelist(*domain),
			}

			// Configurare TLS/SSL a serverului http
			s.TLSConfig = &tls.Config{
				GetCertificate: autoTLSManager.GetCertificate,
				NextProtos:     []string{acme.ALPNProto},
			}

			log.Fatal(s.ListenAndServeTLS("", ""))
		}
	} else {
		// server http
		log.Println("SSL/TLS disabled")
		log.Fatal(s.ListenAndServe())
	}
}

// Intoarce lista de adrese tinta pt proxy
func addTarget(lista []string) (ret []*middleware.ProxyTarget, err error) {

	ret = []*middleware.ProxyTarget{}

	for _, v := range lista {
		var url *url.URL
		url, err = url.Parse(v)
		if err != nil {
			return
		}
		target := &middleware.ProxyTarget{URL: url}
		ret = append(ret, target)
	}

	return
}

// Intoarce lista de tinte
func parseTargets(commaStrIn string) (ret []string, err error) {
	if len(commaStrIn) == 0 {
		err = fmt.Errorf("Proxy target list is empty!")
	}

	ret = strings.Split(commaStrIn, ",")
	log.Printf("%#v\n", ret)
	return
}
