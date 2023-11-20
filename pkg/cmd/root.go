package cmd

import (
	"context"
	"os"

	"github.com/getopends/opends/internal"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg internal.Config

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "opends-server",
		RunE: runCmd(doServe),
	}

	cmd.AddCommand(
		serveCmd(),
		migrateCmd(),
	)

	cmd.PersistentFlags().String("config", "", "Config")
	viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))

	return cmd
}

type runFunc func(ctx context.Context, cmd *cobra.Command, args []string, cfg *internal.Config) error

func runCmd(runner runFunc) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		return runner(cmd.Context(), cmd, args, cfg)
	}
}

func loadConfig() (*internal.Config, error) {
	cfgFile := viper.GetString("config")
	if cfgFile == "" {
		cfgFile = os.Getenv("OPENDS_CONFIG_FILE")
	}

	return internal.NewConfig(cfgFile)
}
