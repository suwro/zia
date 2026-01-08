package config

import (
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
)

// cache allowed ip - for fast running on next request
var allowedIPsCache = map[string]bool{}

// Check if the client IP address exists in Permitted list
// Use memory cache for ip list
func AllowAccessMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			remoteIP := c.RealIP()
			allowed := false
			inCache := false

			// Check valid ip if not exists in application memory cache
			if exists, ok := allowedIPsCache[remoteIP]; ok {
				allowed = exists
				inCache = true
			}

			// validate ip from config Allowed IP
			if !inCache {
				// Check the client IP list
				for _, ipConfig := range CFG.AllowedIPs {
					if validateIP(remoteIP, ipConfig.Ip) {
						allowed = true
						break
					}
				}
				// set the result in mem cache
				allowedIPsCache[remoteIP] = allowed
			}

			// Return default page instead of proxy target pages
			if !allowed {
				return c.Render(http.StatusOK, CFG.NoAccessPage, echo.Map{
					"domainName": CFG.DomainName,
					"ipAddress":  remoteIP,
				})
			}

			// IP is allowed - continue with proxy to target
			return next(c)
		}
	}
}

// Validate if the client IP address exists in Permitted list
// ipConfig can be a simple ip address or a CIDR class
// Example:
//
//	validateIP("10.1.1.1", "10.1.1.0/24") => true
//	validateIP("10.1.1.1","192.168.200.1") => false
func validateIP(ipClient string, ipConfig string) bool {
	// Parsează IP-ul clientului
	clientIP := net.ParseIP(ipClient)
	if clientIP == nil {
		return false
	}

	// Verifică dacă ipConfig este CIDR sau IP simplu
	_, ipNet, err := net.ParseCIDR(ipConfig)
	if err == nil {
		// ipConfig este CIDR - verifică apartenența
		return ipNet.Contains(clientIP)
	}

	// ipConfig este IP simplu - verifică egalitatea
	configIP := net.ParseIP(ipConfig)
	if configIP == nil {
		return false
	}

	return clientIP.Equal(configIP)
}
