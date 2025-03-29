package middlewares

import (
	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenAuthMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		sCookie, err := c.Request.Cookie(session.DefaultCookieName)
		if err != nil {
			if tkn := c.Request.Header.Get("Authorization"); tkn != "" && session.Cache.ValidateToken(tkn) {
				c.Set(gwt.CtxToken, tkn)
				c.Next()
				return
			}
			c.AbortWithError(http.StatusUnauthorized, model.ErrSessionNotFound)
			return
		} else if sCookie.Valid() != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		if !session.Cache.ValidateSession(sCookie.Value) {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		sess := session.Cache.Get(sCookie.Value)
		if sess == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		c.Set(gwt.CtxToken, sess.Token())
		c.Next()
	}
}
