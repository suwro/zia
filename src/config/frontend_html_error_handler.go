package config

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Show error html page, masking internal ip and real error message
func FrontendHtmlErrorHandler(c echo.Context, err error) error {
	log.Println("Error:", err.Error())

	return c.Render(http.StatusGatewayTimeout, "error", echo.Map{
		"domainName": CFG.DomainName,
		"ipAddress":  c.RealIP(),
	})
}
