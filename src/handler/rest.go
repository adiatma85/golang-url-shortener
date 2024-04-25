package handler

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/adiatma85/golang-url-shortener/src/business/usecase"
	"github.com/adiatma85/golang-url-shortener/utils/config"
	"github.com/adiatma85/own-go-sdk/appcontext"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/instrument"
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/adiatma85/own-go-sdk/log"
	"github.com/adiatma85/own-go-sdk/parser"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	infoRequest  string = `httpclient Sent Request: uri=%v method=%v`
	infoResponse string = `httpclient Received Response: uri=%v method=%v resp_code=%v`
)

var once = &sync.Once{}

type REST interface {
	Run()
}

type rest struct {
	http       *gin.Engine
	conf       config.GinConfig
	json       parser.JSONInterface
	log        log.Interface
	uc         *usecase.Usecase
	instrument instrument.Interface
	jwtAuth    jwtAuth.Interface
}

type InitParam struct {
	Http       *gin.Engine
	Conf       config.GinConfig
	Json       parser.JSONInterface
	Log        log.Interface
	Uc         *usecase.Usecase
	Instrument instrument.Interface
	JwtAuth    jwtAuth.Interface
}

func Init(param InitParam) REST {
	r := &rest{}

	once.Do(func() {
		switch param.Conf.Mode {
		case gin.ReleaseMode:
			gin.SetMode(gin.ReleaseMode)
		case gin.DebugMode, gin.TestMode:
			gin.SetMode(gin.TestMode)
		default:
			gin.SetMode("")
		}

		httpServer := gin.New()

		r = &rest{
			conf:       param.Conf,
			log:        param.Log,
			json:       param.Json,
			http:       httpServer,
			uc:         param.Uc,
			instrument: param.Instrument,
			jwtAuth:    param.JwtAuth,
		}

		// Set CORS
		switch r.conf.CORS.Mode {
		case "allowall":
			r.http.Use(cors.New(cors.Config{
				AllowAllOrigins: true,
				AllowHeaders:    []string{"*"},
				AllowMethods: []string{
					http.MethodHead,
					http.MethodGet,
					http.MethodPost,
					http.MethodPut,
					http.MethodPatch,
					http.MethodDelete,
				},
			}))
		default:
			r.http.Use(cors.New(cors.DefaultConfig()))
		}

		// Set Timeout
		r.http.Use(r.SetTimeout)

		// Set Recovery
		r.http.Use(r.CustomRecovery)

		r.Register()
	})

	return r
}

func (r *rest) CustomRecovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// Check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok { // nolint: errorlint
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}
			if brokenPipe {
				// If the connection is dead, we can't write a status to it.
				ctx.Error(err.(error)) // nolint: errcheck
				ctx.Abort()
			} else {
				r.httpRespError(ctx, errors.NewWithCode(codes.CodeInternalServerError, http.StatusText(http.StatusInternalServerError)))
			}

			// Need to update SDK First before uncomment this
			r.log.Panic(err)
		}
	}()
	ctx.Next()
}

func (r *rest) Run() {
	// Create context that listens for the interrupt signal from the OS.
	c := appcontext.SetServiceVersion(context.Background(), r.conf.Meta.Version)
	ctx, stop := signal.NotifyContext(c, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := ":8080"
	if r.conf.Port != "" {
		port = fmt.Sprintf(":%s", r.conf.Port)
	}

	srv := &http.Server{
		Addr:              port,
		Handler:           r.http,
		ReadHeaderTimeout: 2 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.log.Error(ctx, fmt.Sprintf("Serving HTTP error: %s", err.Error()))
		}
	}()
	r.log.Info(ctx, fmt.Sprintf("Listening and Serving HTTP on %s", srv.Addr))

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	r.log.Info(ctx, "Shutting down server...")

	// The context is used to inform the server it has timeout duration to finish
	// the request it is currently handling
	quitctx, cancel := context.WithTimeout(c, r.conf.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(quitctx); err != nil {
		r.log.Fatal(quitctx, fmt.Sprintf("Server Shutdown: %s", err.Error()))
	}
	r.log.Info(quitctx, "Server Shut Down.")
}

func (r *rest) Register() {
	// server health and testing purpose
	r.http.GET("/ping", r.Ping)
	// r.registerSwaggerRoutes()
	// r.registerDummyRoutes()

	// // Set Common Middlewares
	// commonPublicMiddlewares := gin.HandlersChain{
	// 	r.addFieldsToContext, r.BodyLogger,
	// }

	// commonPrivateMiddlewares := gin.HandlersChain{
	// 	r.addFieldsToContext, r.BodyLogger,
	// 	r.VerifyUser,
	// }

	// // public api
	// publicv1 := r.http.Group("/public/v1/", commonPublicMiddlewares...)
	// publicv1.POST("/register", r.RegisterNewUserWithoutToken)

	// // auth api
	// authv1 := r.http.Group("/auth/v1", commonPublicMiddlewares...)
	// authv1.POST("/login", r.SignInWithPassword)
	// authv1.GET("/refresh-token", r.VerifyUser, r.RefreshToken)

	// // private api
	// v1 := r.http.Group("/v1/", commonPrivateMiddlewares...)

	// // user
	// v1.GET("/user/:user_id", r.GetUserByID)
	// v1.GET("/user/profile", r.UserProfile)
	// v1.PUT("/user/profile", r.UpdateUserProfile)
	// v1.DELETE("/user/profile", r.UserSelfDelete)
	// v1.PUT("/user/profile/change-password", r.UserChangePassword)

	// // user management admin api
	// v1.GET("/admin/user", r.isAdmin, r.GetListUserAsAdmin)
	// v1.DELETE("/admin/user/:user_id", r.DeleteUser)
	// v1.PUT("/admin/user/:user_id", r.isAdmin, r.UpdateUser)

	// // category
	// v1.GET("/category", r.GetListCategory)
	// v1.POST("/category", r.CreateCategory)
	// v1.GET("/category/:category_id", r.GetCategoryByID)
	// v1.PUT("/category/:category_id", r.UpdateCategory)
	// v1.DELETE("/category/:category_id", r.DeleteCategory)

	// // task
	// v1.GET("/task", r.GetListTask)
	// v1.POST("/task", r.CreateTask)
	// v1.GET("/task/:task_id", r.GetTaskById)
	// v1.PUT("/task/:task_id", r.UpdateTask)
	// v1.DELETE("/task/:task_id", r.DeleteTask)

	// // role
	// v1.GET("/role", r.isAdmin, r.GetListRole)
	// v1.POST("/role", r.isAdmin, r.CreateRole)
	// v1.GET("/role/:role_id", r.isAdmin, r.GetRoleById)
	// v1.PUT("/role/:role_id", r.isAdmin, r.UpdateRole)
	// v1.DELETE("/role/:role_id", r.isAdmin, r.DeleteRole)
}

// func (r *rest) registerSwaggerRoutes() {
// 	if r.conf.Swagger.Enabled {
// 		swagger.SwaggerInfo.Title = r.conf.Meta.Title
// 		swagger.SwaggerInfo.Description = r.conf.Meta.Description
// 		swagger.SwaggerInfo.Version = r.conf.Meta.Version
// 		swagger.SwaggerInfo.Host = r.conf.Meta.Host
// 		swagger.SwaggerInfo.BasePath = r.conf.Meta.BasePath

// 		swaggerAuth := gin.Accounts{
// 			r.conf.Swagger.BasicAuth.Username: r.conf.Swagger.BasicAuth.Password,
// 		}

// 		r.http.GET(fmt.Sprintf("%s/*any", r.conf.Swagger.Path),
// 			gin.BasicAuthForRealm(swaggerAuth, "Restricted"),
// 			// ginSwagger.WrapHandler(swaggerfiles.Handler))
// 	}
// }

// func (r *rest) registerDummyRoutes() {
// 	if r.conf.Dummy.Enabled {
// 		// load login page to gin

// 		r.http.LoadHTMLFiles(
// 			"./docs/templates/login.html",
// 		)

// 		dummyGroup := r.http.Group(r.conf.Dummy.Path)
// 		fmt.Println(dummyGroup)
// 		// dummyGroup.GET("/login", r.DummyLogin)
// 	}
// }
