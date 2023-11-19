package cmd

import (
	"github.com/getopends/opends/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg internal.Config

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "opends-server",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := internal.NewConfig(viper.GetString("config"))
			if err != nil {
				return err
			}

			return runServe(cfg)
		},
	}

	cmd.AddCommand(
		cmdServe(&cfg),
		cmdMigrate(),
	)

	cmd.PersistentFlags().String("config", "", "Config")
	viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))

	return cmd
}
