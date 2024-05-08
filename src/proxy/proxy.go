package proxy

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var hostProxy map[string]*httputil.ReverseProxy = map[string]*httputil.ReverseProxy{}

type BaseHandle struct{}

func (h *BaseHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.Host, ":")[0]
	log.Println(Au.Yellow(r.RemoteAddr), "host:", Au.Green(host))

	if fn, ok := hostProxy[host]; ok {
		fn.ServeHTTP(w, r)
		return
	}

	if target, ok := hostTarget[host]; ok {
		remoteUrl, err := url.Parse(target)
		if err != nil {
			log.Println(Au.Red("Target parse fail:"), err)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		hostProxy[host] = proxy
		proxy.ServeHTTP(w, r)
		return
	}

	// Randare pagina principala
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmplData := map[string]interface{}{
		"reqHost":  r.Host,
		"hostName": Cfg.Service.HostName,
		"ver":      Version,
	}
	err := tmpl.Execute(w, tmplData)
	if err != nil {
		fmt.Fprint(w, err)
	}

	/*
		_, err := w.Write([]byte("403: Host forbidden " + host))
		if err != nil {
			log.Println(err.Error())
		}
	*/
}
