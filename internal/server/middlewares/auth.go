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
			c.Redirect(http.StatusFound, model.PathLogin)
			c.Abort()
			return
		} else if sCookie.Valid() != nil {
			c.Redirect(http.StatusFound, model.PathLogin)
			c.Abort()
			return
		}
		if sess, ok := session.Cache.ValidateSession(sCookie.Value); !ok || sess == nil {
			c.Redirect(http.StatusFound, model.PathLogin)
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

// NonAuthMw assures that further is guest only
func NonAuthMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool(session.CtxLoggedIn) {
			refs := c.Request.Header[model.RefererHeader]
			if len(refs) == 0 || refs[0] == "" {
				refs = []string{model.PathMain}
			}
			c.Redirect(http.StatusNotModified, refs[0])
			c.Abort()
			return
		}
		c.Next()
	}
}

func AdminAuthMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		sCookie, err := c.Request.Cookie(session.DefaultCookieName)
		if err != nil {
			c.Redirect(http.StatusFound, model.PathLogin)
			c.Abort()
			return
		} else if sCookie.Valid() != nil {
			c.Redirect(http.StatusFound, model.PathLogin)
			c.Abort()
			return
		}
		sess, ok := session.Cache.ValidateSession(sCookie.Value)
		if !ok || sess == nil {
			c.Redirect(http.StatusFound, model.PathLogin)
			c.Abort()
			return
		}

		if tkn := session.Cache.GetJWTBySessionID(sess.ID); tkn != nil && tkn.Role == model.RoleAdmin {
			c.Next()
			return
		}
		c.Redirect(http.StatusFound, model.PathLogin)
		c.Abort()
		return
	}
}
