package config

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Verifica daca adresa IP este in lista de IP-uri permise
func AllowAccessMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extrage adresa IP a cererii
			remoteIP := c.RealIP()

			// Verifica daca adresa IP este in lista de IP-uri permise
			allowed := false
			for _, ipConfig := range CFG.AllowedIPs {
				if ipConfig.Ip == remoteIP {
					allowed = true
					break
				}
			}

			// Daca adresa IP nu este permisă, intoarce mesajul "OK"
			if !allowed {
				return c.Render(http.StatusOK, CFG.NoAccessPage, echo.Map{
					"domainName": CFG.DomainName,
					"ipAddress":  remoteIP,
				})
			}

			// Daca adresa IP este permisă, continuează cu următorul handler
			return next(c)
		}
	}
}
