package config

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Verifica daca adresa IP este in lista de IP-uri permise
func AllowAccessMiddleware(allowedIPs []AllowedIP, domainName, noAccessPage string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extrage adresa IP a cererii
			remoteIP := c.RealIP()

			println("DomainName:", domainName)

			// Verifica daca adresa IP este in lista de IP-uri permise
			allowed := false
			for _, ipConfig := range allowedIPs {
				if ipConfig.Ip == remoteIP {
					allowed = true
					break
				}
			}

			// Daca adresa IP nu este permisă, intoarce mesajul "OK"
			if !allowed {
				return c.Render(http.StatusOK, noAccessPage, echo.Map{
					"domainName": domainName,
					"ipAddress":  remoteIP,
				})
			}

			// Daca adresa IP este permisă, continuează cu următorul handler
			return next(c)
		}
	}
}
