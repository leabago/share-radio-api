package http_api

import (
	"net/http"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/leabago/share-radio/adder/config"
	"github.com/leabago/share-radio/adder/docs"
	"github.com/leabago/share-radio/adder/docs/gen"
	"github.com/leabago/share-radio/adder/internal/controller/http_api/api"
	"github.com/leabago/share-radio/adder/internal/controller/http_api/middleware"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/leabago/share-radio/adder/internal/usecase"
	"github.com/leabago/share-radio/adder/pkg/logger"
)

func NewRouter(app *fiber.App, cfg *config.Config,
	station usecase.Station, genre usecase.Genre, language usecase.Language, l logger.Interface) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))
	app.Use(cors.New(cors.Config{

		AllowOrigins: "*",
	}))

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New("my-service-name")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	// Swagger
	if cfg.Swagger.Enabled {

		app.Use(swagger.New(swagger.Config{
			BasePath:    "/",
			FileContent: docs.ApiSchema,
			Path:        "swagger",

			CacheAge: 1,
		}))

	}

	// K8s probe
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	// Routers
	bffController := api.NewController(station, genre, language)
	wrappedHandler := gen.NewStrictHandler(bffController, nil)
	gen.RegisterHandlers(app, wrappedHandler)

}
