package bootstrap

import (
	"go-authorization/errors"
	"go-authorization/internal/controllers"
	"go-authorization/internal/middlewares"
	"go-authorization/internal/repository"
	"go-authorization/internal/routes"
	services "go-authorization/internal/service"
	"go-authorization/lib"
	"net/http"
	"time"
)

type Application struct {
	Handler     lib.HttpHandler
	Database    lib.Database
	Logger      lib.Logger
	Config      lib.Config
	Middleware  middlewares.Middlewares
	Redis       lib.Redis
	Repository  repository.Module
	Services    services.Module
	Controllers controllers.Module
	Routes      routes.Module
}

func New() {
	app := &Application{}
	app.Config = lib.NewConfig()
	app.Logger = lib.NewLogger(app.Config)
	app.Handler = lib.NewHttpHandler(app.Logger, app.Config)
	app.Database = lib.NewDatabase(app.Config, app.Logger)
	app.Redis = lib.NewRedis(app.Config, app.Logger)
	app.ProvideRepository()
	app.ProvideService()
	app.ProvideControllers()

	app.ProvideMiddleware()
	app.ProvideRoutes()

	app.Start()
}

func (app *Application) Start() {
	db, err := app.Database.ORM.DB()
	if err != nil {
		app.Logger.Zap.Fatalf("Error to get database connection: %v", err)
	}

	app.Logger.Zap.Info("Starting Application")

	if err := db.Ping(); err != nil {
		app.Logger.Zap.Fatalf("Error to ping database connection: %v", err)
	}
	// set conn
	db.SetMaxOpenConns(app.Config.Database.MaxOpenConns)
	db.SetMaxIdleConns(app.Config.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(app.Config.Database.MaxLifetime) * time.Second)

	// Setting Middleware
	app.Middleware.Setup()
	app.Routes.Setup(app.Logger)

	//e.GET("/metrics", echoprometheus.NewHandler())
	//e.GET("/healthcheck", func(c echo.Context)

	// Engine Start
	if err := app.Handler.Engine.Start(app.Config.Http.ListenAddr()); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			app.Logger.Zap.Debug("Shutting down the Application")
		} else {
			app.Logger.Zap.Fatalf("Error to Start Application: %v", err)
		}
	}
}

func (app *Application) ProvideMiddleware() {
	app.Middleware = middlewares.NewMiddlewares(
		middlewares.NewCoreMiddleware(app.Handler, app.Logger, app.Database),
		middlewares.NewCorsMiddleware(app.Handler, app.Logger),
		middlewares.NewZapMiddleware(app.Handler, app.Logger),
		middlewares.NewAuthMiddleware(app.Handler, app.Logger, app.Config, services.NewAuthService(app.Redis, app.Config)),
		middlewares.NewCasbinMiddleware(app.Handler, app.Logger, app.Config, app.Services.CasbinService))
}
