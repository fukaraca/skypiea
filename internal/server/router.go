package service

import (
	"github.com/fukaraca/skypiea/internal/config"
	"github.com/fukaraca/skypiea/internal/service/templater"
	logg "github.com/fukaraca/skypiea/pkg/log"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func NewRouter(cfg *config.Server, logger *slog.Logger, opts ...gin.OptionFunc) *gin.Engine {
	gin.SetMode(cfg.GinMode)
	e := gin.New(opts...)
	templates := templater.New()
	templates.LoadHTMLGlob("./web/templates")
	e.HTMLRender = templates
	e.Use(commonMiddlewares(logger)...)
	return e
}

func commonMiddlewares(logger *slog.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.Recovery(),
		static.Serve("/web/static", static.LocalFile("./web/static", false)),
		logg.GinMiddleware(logger.With("via", "rest")),
	}
}

func (s *Server) bindRoutes() {
	s.engine.NoRoute()

	s.engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{
			"Title":    "Home",
			"CSSFile":  "index.css",
			"LoggedIn": false,
		})
	})
	s.engine.GET("/pricing", func(c *gin.Context) {
		c.HTML(http.StatusOK, "pricing", gin.H{
			"Title": "Pricing",
		})
	})
	s.engine.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login", gin.H{
			"Title": "Login",
		})
	})
	s.engine.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup", gin.H{
			"Title": "Sign Up",
		})
	})
	s.engine.GET("/forgot-password", func(c *gin.Context) {
		c.HTML(http.StatusOK, "forgot-password", gin.H{
			"Title": "Recover your password",
		})
	})
	s.engine.GET("/profile", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile", gin.H{
			"Title": "My Profile",
		})
	})
}
