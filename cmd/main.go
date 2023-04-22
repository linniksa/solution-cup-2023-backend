package main

import (
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"
	"solution/config"
	"solution/controllers"
	"solution/ent"
	"solution/rates"
)

func main() {
	unix.Umask(0)

	app := fx.New(
		fx.Provide(NewLogger),

		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger.WithOptions(zap.AddCallerSkip(1))}
		}),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("solution")
		}),

		fx.Provide(
			config.New,
			ent.New,
			rates.New,
		),

		controllers.Fx(),
	)

	app.Run()
}

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.Dev {
		return zap.NewDevelopment()
	}

	return zap.NewProduction()
}
