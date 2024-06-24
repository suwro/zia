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

var versiune = "0.3.1"

func main() {
	domain := flag.String("domain", "", "Domain name to use for the certificate")
	port := flag.Int("port", 8080, "Port to listen on")
	ssl := flag.Bool("ssl", true, "Use SSL/TLS default true")
	ver := flag.Bool("version", false, "Show version and exit")
	targetList := flag.String("targets", "", "List of targets for proxy, comma separated")
	timeout := flag.Int("timeout", 0, "Timeout for proxy in seconds, 0 no timeout")

	flag.Parse()
	log.Println("Versiune:", versiune)
	if *ver {
		os.Exit(1)
	}

	// verifica domeniu
	if len(*domain) == 0 {
		log.Fatal("Domain name is required")
	}

	// verifica lista de tinte proxy
	lista, err := parseTargets(*targetList)
	if err != nil {
		log.Fatal(err)
	}

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

	//c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXMLCharsetUTF8)

	// Server http proxy
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}\t${remote_ip}\t${method}\t${uri}\t${status} ${error}\n",
		Output: f,
	}))

	// TLS Transport proxy
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
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
			MinVersion:     tls.VersionTLS13,
		}

		log.Fatal(s.ListenAndServeTLS("", ""))
	} else {
		// server http
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
	return
}
