package middlewares

import (
	"net/http"

	"github.com/fukaraca/skypiea/internal/model"
	"github.com/fukaraca/skypiea/pkg/gwt"
	"github.com/fukaraca/skypiea/pkg/session"
	"github.com/gin-gonic/gin"
)

func StrictAuthMw() gin.HandlerFunc {
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
			c.Header(model.HxRedirect, "/login")
			c.Status(http.StatusFound)
			c.Abort()
			return
		} else if sCookie.Valid() != nil {
			c.Header(model.HxRedirect, "/login")
			c.Status(http.StatusFound)
			c.Abort()
			return
		}
		if sess, ok := session.Cache.ValidateSession(sCookie.Value); !ok || sess == nil {
			c.Header(model.HxRedirect, "/login")
			c.Status(http.StatusFound)
			c.Abort()
			return
		} else {
			c.Set(gwt.CtxToken, sess.Token())
			c.Next()
			return
		}
	}
}

func CommonAuthMw() gin.HandlerFunc {
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

func NonAuthMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(session.CtxLoggedIn) {
			refs := c.Request.Header[model.RefererHeader]
			if len(refs) == 0 || refs[0] == "" {
				refs = []string{"/"}
			}
			c.Header(model.HxRedirect, refs[0])
			c.Status(http.StatusNotModified)
			c.Abort()
			return
		}
		c.Next()
	}
}
