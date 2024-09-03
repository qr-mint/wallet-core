package api

import (
	"errors"
	"fmt"
	_ "github.com/didip/tollbooth"
	_ "github.com/didip/tollbooth/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"nexus-wallet/api/controllers"
	"nexus-wallet/api/error_handler"
	"nexus-wallet/internal/app_error"
	"runtime/debug"
	"strings"
	"time"
)

type Kernel struct {
	routers      []controllers.Router
	apiLogger    io.Writer
	errorHandler *error_handler.HttpErrorHandler
	params       Params
}

type Params struct {
	Port            string
	HeaderOriginURL string
	DebugMode       bool
}

func NewKernel(
	routers []controllers.Router,
	apiLogger io.Writer,
	errorHandler *error_handler.HttpErrorHandler,
	params Params,
) Kernel {
	return Kernel{
		routers:      routers,
		apiLogger:    apiLogger,
		errorHandler: errorHandler,
		params:       params,
	}
}

func (s *Kernel) Run() error {
	gin.Default()
	if !s.params.DebugMode {
		gin.SetMode(gin.ReleaseMode)
		println("started http kernel")
	}
	engine := gin.New()

	s.registerLogger(engine)
	s.registerRecovery(engine)
	s.registerStatic(engine)
	s.registerCors(engine)
	s.register404(engine)
	s.registerRoutes(engine)

	return engine.Run(":" + s.params.Port)
}

func (s *Kernel) register404(engine *gin.Engine) {
	engine.NoRoute(func(context *gin.Context) {
		s.errorHandler.Handle(context, app_error.ResourceNotFoundError(errors.New("route not found")))
		context.Next()
	})
}

func (s *Kernel) registerRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")
	{
		for _, route := range s.routers {
			route.SetRoutes(v1)
		}
	}
}

func (s *Kernel) registerRecovery(engine *gin.Engine) {
	engine.Use(func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				s.errorHandler.Handle(
					context,
					app_error.InternalError(fmt.Errorf("panic: %s. stack: %s", err, debug.Stack())),
				)
			}
		}()
		context.Next()
	})
}

func (s *Kernel) registerLogger(engine *gin.Engine) {
	engine.Use(gin.LoggerWithWriter(s.apiLogger))
}

func (s *Kernel) registerStatic(engine *gin.Engine) {
	engine.Static("/images", "./assets/images")
}

func (s *Kernel) registerCors(engine *gin.Engine) {
	headerString := `host,connection,sec-ch-ua,sec-ch-ua-mobile,sec-ch-ua-platform,upgrade-insecure-requests,user-agent,accept,sec-fetch-site,sec-fetch-mode,sec-fetch-user,sec-fetch-dest,accept-encoding,accept-language,cookie
	content-type,host,connection,sec-ch-ua,sec-ch-ua-mobile,user-agent,sec-ch-ua-platform,accept,sec-fetch-site,sec-fetch-mode,sec-fetch-dest,referer,accept-encoding,accept-language,cookie`

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{s.params.HeaderOriginURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    strings.Split(headerString, ","),
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
