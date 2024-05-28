package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/logrusorgru/aurora"
)

var (
	Version    string
	Cfg        *ConfigStruct
	hostTarget = map[string]string{}
	Au         aurora.Aurora
)

func init() {
	Cfg = &ConfigStruct{}
	Au = aurora.NewAurora(true)
}

func InitConfig(configFile string) {

	// load config file
	content, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	// parse config file
	//m := map[string]interface{}{}
	err = json.Unmarshal(content, &Cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Validate Config
	validate := validator.New()
	err = validate.Struct(Cfg.Service)
	if err != nil {
		log.Fatal(err.Error())
	}

	// service IP 0.0.0.0 - to show on start
	if len(Cfg.Service.IP) == 0 {
		Cfg.Service.IP = "0.0.0.0"
	}

	// Ck config domains
	if len(Cfg.Domains) == 0 {
		log.Fatal(Au.Red("Config has no domains defined to proxy!"))
	}

	// Validate Domains
	for n, domeniu := range Cfg.Domains {
		nr := fmt.Sprintf("%02d", n+1)
		fmt.Println(nr, domeniu.Host, "=>", domeniu.Target)

		err = validate.Struct(domeniu)
		if err != nil {
			log.Fatal(err.Error())
		}
		hostTarget[domeniu.Host] = domeniu.Target
	}
}
