package cmd

import (
	"os"

	"github.com/getopends/opends/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg internal.Config

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "opends-server",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfgFile := viper.GetString("config")
			if cfgFile == "" {
				cfgFile = os.Getenv("OPENDS_CONFIG_FILE")
			}

			cfg, err := internal.NewConfig(cfgFile)
			if err != nil {
				return err
			}

			return doServe(cfg)
		},
	}

	cmd.AddCommand(
		serveCmd(&cfg),
		migrateCmd(),
	)

	cmd.PersistentFlags().String("config", "", "Config")
	viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))

	return cmd
}
