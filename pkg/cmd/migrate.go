package cmd

import "github.com/spf13/cobra"

func migrateCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:  "migrate",
		RunE: doMigrate,
	}

	return newCmd
}

func doMigrate(cmd *cobra.Command, args []string) error {
	return nil
}
