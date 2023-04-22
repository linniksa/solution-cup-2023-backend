package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/lock ./schema

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect"
	"fmt"
	"go.uber.org/fx"
	config2 "solution/config"
)

import entSQL "entgo.io/ent/dialect/sql"

func (c *Client) DB() *sql.DB {
	if c.debug {
		return c.driver.(*dialect.DebugDriver).Driver.(*entSQL.Driver).DB()
	}

	return c.driver.(*entSQL.Driver).DB()
}

func New(
	lc fx.Lifecycle,
	cfg *config2.Config,
) (*Client, error) {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	client, err := Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.StartStopHook(
		func(ctx context.Context) error {
			// Run the auto migration tool.
			if err := client.Schema.Create(ctx); err != nil {
				return fmt.Errorf("failed creating schema resources: %w", err)
			}

			return nil
		},
		func(ctx context.Context) error {
			return client.Close()
		},
	))

	return client, nil
}
