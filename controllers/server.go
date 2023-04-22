package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"solution/config"
	"solution/ent"
)

var controllers []interface{}

func Fx() fx.Option {
	return fx.Module("api",
		fx.Provide(
			NewHttpServer,
			func(app *fiber.App) fiber.Router {
				return app.Group("/api/v1")
			},
		),
		fx.Invoke(controllers...),
	)
}

func NewHttpServer(
	lc fx.Lifecycle,
	logger *zap.Logger,
	appCfg *config.Config,
) (*fiber.App, error) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}

	fiberCfg := fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			{
				var e *fiber.Error
				if errors.As(err, &e) {
					code = e.Code
				}
			}

			{
				var e *ent.NotFoundError
				if errors.As(err, &e) {
					code = 400
				}
			}

			if err != nil && (code >= 500 || appCfg.Dev) {
				logger.Error(
					"http server err",
					zap.Error(err),
					zap.String("path", c.Path()),
				)
			}

			// Set Content-Type: text/plain; charset=utf-8
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

			if code >= 500 {
				code = 500
			} else if code >= 400 {
				code = 400
			}

			// Return status code with error message
			return c.Status(code).SendString(err.Error())
		},
	}

	if appCfg.Dev {
		fiberCfg.JSONEncoder = func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", "  ")
		}
	}

	app := fiber.New(fiberCfg)
	lc.Append(fx.StartHook(func() error {
		go app.Listener(ln)

		return nil
	}))

	app.Use(recover.New())
	// app.Use(compress.New())

	if appCfg.Dev {
		app.Use(func(c *fiber.Ctx) (err error) {

			logger.Info("new request",
				zap.String("req", c.Request().String()),
				zap.ByteString("body", c.Request().Body()),
			)

			return c.Next()
		})
	}

	return app, nil
}
