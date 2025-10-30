package cmd

import (
	"context"
	"fmt"

	"github.com/kborovik/gcp-iam/config"
	"github.com/kborovik/gcp-iam/db"
	"github.com/urfave/cli/v3"
)

// DBAction is a command action that receives config and database instances
type DBAction func(ctx context.Context, cmd *cli.Command, cfg *config.Config, database *db.DB) error

// WithDB wraps a command action with config loading and database initialization.
// This eliminates the need to repeat config and database setup in every command.
func WithDB(action DBAction) func(context.Context, *cli.Command) error {
	return func(ctx context.Context, cmd *cli.Command) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		database, err := db.New(cfg.DatabasePath)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer database.Close()

		return action(ctx, cmd, cfg, database)
	}
}
