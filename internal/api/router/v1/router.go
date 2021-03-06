package v1

import (
	"fmt"
	"net/http"

	"github.com/adiatma85/golang-rest-template-api/internal/api/handler"
	"github.com/adiatma85/golang-rest-template-api/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// V1 Router
func Setup() *gin.Engine {
	app := gin.New()

	// Middlewares
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - - [%s] \"%s %s %s %d %s \" \" %s\" \" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	app.Use(gin.Recovery())
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	// Routes for v1
	v1Route := app.Group("/api/v1")

	// Home to return what status it is
	{
		v1Route.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		})
	}

	// AuthGroup with "auth" prefix
	authGroup := v1Route.Group("auth")
	authHandler := handler.GetAuthHandler()
	{
		authGroup.POST("login", authHandler.Login)
		authGroup.POST("register", authHandler.Register)
		profileAuthGroup := authGroup.Group("profile")
		profileAuthGroup.Use(middleware.AuthJWT(), middleware.IsAdminOrUserMiddleware())
		{
			profileAuthGroup.GET("", authHandler.GetProfile)
			profileAuthGroup.PUT("", authHandler.UpdateProfile)
			profileAuthGroup.DELETE("", authHandler.DeleteProfile)
		}
	}

	// UrlGroup with "url" prefix
	urlGroup := v1Route.Group("url")
	urlHandler := handler.GetUrlHandler()
	urlGroup.Use(middleware.AuthJWT())
	{
		urlGroup.POST("", middleware.IsAdminOrUserMiddleware(), urlHandler.Create)
		urlGroup.GET("query", middleware.IsAdminOrUserMiddleware(), urlHandler.Query)
		urlGroup.GET("load/:short_token", middleware.IsAdminOrUserMiddleware(), urlHandler.Load)

		// Authorized Url Group with "url/authorized" prefix
		authorizedUrlGroup := urlGroup.Group("authorized")
		authorizedUrlGroup.Use(middleware.IsAdminMiddleware())
		{
			authorizedUrlGroup.POST("", urlHandler.AuthorizedCreate)
			authorizedUrlGroup.PUT(":id", urlHandler.AuthorizedUpdate)
			authorizedUrlGroup.DELETE(":id", urlHandler.AuthorizedDelete)
		}
	}
	return app
}
