package cmd

import "github.com/spf13/cobra"

func cmdMigrate() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use: "migrate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return migrateCmd
}
