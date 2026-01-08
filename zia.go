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
	"time"
	"zia/src/config"

	"golang.org/x/crypto/acme"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
)

// Versiunea aplicatiei
var versiune = "0.3.5"

func main() {
	dev := flag.Bool("dev", false, "Dev mode")
	stdout := flag.Bool("stdout", false, "Use stdout instead /var/log/zia/<domain>/acces_<port>.log")
	ver := flag.Bool("version", false, "Show version and exit")
	configfile := flag.String("config", "", "Configuration file (not implemented yet)")

	flag.Parse()
	log.Println("Versiune:", versiune)
	if *ver {
		os.Exit(1)
	}

	// Pointer to config
	cfg := &config.CFG

	// Config File
	if len(*configfile) > 0 {
		log.Println("Reading configuration from:", *configfile)
		config.Parse(*configfile)
	}

	// Dev Mode
	if *dev {
		log.Println("Dev mode enabled")
		config.DisableCache = true
		*stdout = true // default stdout in devmode
	}

	// Logs
	var logConfig = middleware.LoggerConfig{
		Format: "${time_rfc3339}\t${remote_ip}\t${method}\t${uri}\t${status} ${error}\n",
	}

	// Setare log-uri in fisier in loc de consola
	if !*stdout {
		// log in fisier in loc de standard
		logFileName := fmt.Sprintf("/var/log/zia/%s/acces_%d.log", cfg.DomainName, cfg.Port)

		// Salveaza log-urile in logs.txt
		f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Error generating log file: %v", err)
		}

		// activeaza magia log-urilor
		defer func(file *os.File) {
			if err = file.Close(); err != nil {
				log.Fatal(err.Error())
			}
		}(f)

		log.SetOutput(f)
		logConfig.Output = f
	}

	// Server http proxy
	e := echo.New()
	e.Renderer = config.SetRenderer()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(logConfig))

	// Middleware Access List
	if len(cfg.AllowedIPs) > 0 {
		e.Use(config.AllowAccessMiddleware()) //cfg.AllowedIPs, cfg.DomainName, cfg.NoAccessPage
	}

	// TLS Transport proxy
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			//MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		},
	}

	// Proxy
	targets, err := addTarget(cfg.Targets)
	if err != nil {
		log.Fatal(err.Error())
	}

	balancer := middleware.NewRoundRobinBalancer(targets)
	e.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer:     balancer,
		Transport:    transport,
		ErrorHandler: config.FrontendHtmlErrorHandler,
	}))

	// http settings
	s := http.Server{
		Addr:    fmt.Sprintf(":%d", config.CFG.Port),
		Handler: e, // Echo instance handler
	}

	// timeout
	if (cfg.Timeout > 0) && (cfg.Timeout < 3600) {
		s.ReadTimeout = time.Duration(cfg.Timeout) * time.Second
	} else if cfg.Timeout > 3600 {
		log.Fatal("Timeout too high, set it to 0 for no timeout")
	}

	// Show IP/Port message
	log.Println("Prepparing Server Ip:", cfg.Ip, "Port:", cfg.Port)

	// server https
	if cfg.SSL {
		log.Println("SSL/TLS enabled")

		s.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}

		// Given cert and key
		if len(cfg.Cert) > 0 || len(cfg.Key) > 0 {
			log.Fatal(s.ListenAndServeTLS(cfg.Cert, cfg.Key))
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
				HostPolicy: autocert.HostWhitelist(cfg.DomainName),
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
