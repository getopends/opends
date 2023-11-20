package cmd

import (
	"context"

	"github.com/getopends/opends/internal"
	"github.com/spf13/cobra"
)

func migrateCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:  "migrate",
		RunE: runCmd(doMigrate),
	}

	return newCmd
}

func doMigrate(ctx context.Context, cmd *cobra.Command, _ []string, cfg *internal.Config) error {
	cmd.Println("migrate")

	return nil
}
