package config

import (
	"html/template"
	"time"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/echoview-v4"
)

var (
	Root            = "static/template"
	Master          = "master.html"
	DisableCache    = false
	RenderPartials  = []string{} //[]string{"partials/menu"}
	RenderFunctions = template.FuncMap{
		"copy": func() string {
			return time.Now().Format("2006")
		},
		"safeHTML": func(v string) template.HTML {
			return template.HTML(v)
		},
		/*
			"versiuneHash": func() string {
				vh := md5.Sum([]byte(Versiune))
				return hex.EncodeToString(vh[:])
			},
		*/
	}
)

/*
// Adauga render functions la template
	config.RenderFunctions["safeHtml"]= func(v string) template.HTML {
		return template.HTML(v)
	}
*/

func SetRenderer() *echoview.ViewEngine {
	gCfg := goview.DefaultConfig
	gCfg.Root = Root
	gCfg.Master = Master
	gCfg.DisableCache = DisableCache
	gCfg.Partials = RenderPartials

	return echoview.New(gCfg)
}
