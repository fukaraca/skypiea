package middlewares

import (
	"net"
	"net/http"
	"strings"

	gb "github.com/fukaraca/skypiea/pkg/guest_book"
	"github.com/gin-gonic/gin"
)

func CounterUIMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		rip := realIP(c.Request)
		gb.GuestBook.RegisterGuest(rip, c.Request.URL.Path)
		c.Next()
	}
}

func realIP(r *http.Request) string {
	if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
		parts := strings.Split(xf, ",")
		return strings.TrimSpace(parts[0])
	}
	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return xr
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}
