package middlewares

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func TokenAuthMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		sCookie, err := c.Request.Cookie(session.DefaultCookieName)
		if err != nil {
			if tkn := c.Request.Header.Get("Authorization"); tkn != "" {
				if session.Cache.ValidateToken(tkn) {
					c.Set(gwt.CtxToken, tkn)
					c.Next()
					return
				}
				c.AbortWithError(http.StatusUnauthorized, model.ErrInvalidToken)
				return
			}
			c.AbortWithError(http.StatusUnauthorized, model.ErrSessionNotFound)
			return
		} else if sCookie.Valid() != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		if sess, ok := session.Cache.ValidateSession(sCookie.Value); !ok || sess == nil {
			c.Error(model.ErrSessionNotValid)
			c.Redirect(http.StatusFound, "/login")
			return
		} else {
			c.Set(gwt.CtxToken, sess.Token())
			c.Next()
			return
		}
	}
}

func ViewTokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sCookie, err := c.Request.Cookie(session.DefaultCookieName)
		if err != nil || sCookie.Valid() != nil {
			c.Set(session.CtxLoggedIn, false)
			c.Next()
			return
		}
		if sess, ok := session.Cache.ValidateSession(sCookie.Value); !ok || sess == nil {
			c.Set(session.CtxLoggedIn, false)
			c.Next()
			return
		} else {
			c.Set(session.CtxLoggedIn, true)
			c.Set(gwt.CtxToken, sess.Token())
			c.Next()
			return
		}
	}
}
