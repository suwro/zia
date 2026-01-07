package config

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
)

var CFG Config

// Proceseaza fisierul de configurare
func Parse(fileName string) {
	bFileContent, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Proceseaza fisierul de configurare
	err = json.Unmarshal(bFileContent, &CFG)
	if err != nil {
		log.Fatal(err.Error())
	}

	// for Debug
	cfg := &CFG

	// Verifica daca exista toate sectiunile necesare
	if len(cfg.Ip) == 0 {
		cfg.Ip = "0.0.0.0"
	}
	if cfg.Port == 0 {
		cfg.Port = 443
	}
	if len(cfg.DomainName) == 0 {
		log.Fatal("Config file error: DomainName is missing")
	}
	// Validate Targets
	for _, target := range cfg.Targets {
		if !validateUrl(target) {
			log.Fatal("Config file error: Invalid Targets:", target)
		}
	}
}

// Validare Format URL si adresa IP din URL
func validateUrl(url string) bool {
	// validate url - first format http https://...
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}

	// Extract IP address from url
	_tmp := strings.Split(url, ":")
	if len(_tmp) < 2 {
		return false
	}

	// fmt.Println("_tmt: %#v, %d", _tmp, len(_tmp))

	// remove // from ip
	ipAddress := _tmp[1][2:]
	if len(ipAddress) == 0 {
		return false
	}

	// Validate IP address
	if net.ParseIP(ipAddress) != nil {
		return true
	}

	return false
}
