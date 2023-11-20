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
		Use:          "opends-server",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		serveCmd(),
		migrateCmd(),
	)

	registerConfigFlags(cmd)

	return cmd
}

func registerConfigFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("config", "", "Config")
	viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))
}

func registerServeFlags(cmd *cobra.Command) {
	cmd.Flags().String("public-host", "", "Public host")
	viper.BindPFlag("public.host", cmd.Flags().Lookup("public-host"))

	cmd.Flags().String("public-port", "", "Public port")
	viper.BindPFlag("public.port", cmd.Flags().Lookup("public-port"))

	cmd.Flags().String("db-driver", "", "Database driver")
	viper.BindPFlag("database.driver", cmd.Flags().Lookup("db-driver"))

	cmd.Flags().String("db-dsn", "", "Database DSN")
	viper.BindPFlag("database.driver", cmd.Flags().Lookup("db-dsn"))

	cmd.Flags().String("tls-key-file", "", "TLS key file")
	viper.BindPFlag("tls.key_file", cmd.Flags().Lookup("tls-key-file"))

	cmd.Flags().String("tls-cert-file", "", "TLS cert file")
	viper.BindPFlag("tls.cert_file", cmd.Flags().Lookup("tls-cert-file"))

	cmd.Flags().String("tls-key", "", "TLS key")
	viper.BindPFlag("tls.key", cmd.Flags().Lookup("tls-key"))

	cmd.Flags().String("tls-cert", "", "TLS cert")
	viper.BindPFlag("tls.cert", cmd.Flags().Lookup("tls-cert"))
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
