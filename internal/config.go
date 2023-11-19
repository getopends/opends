package internal

import (
	"github.com/spf13/viper"
)

var (
	envPrefix         = "OPENDS"
	defaultConfigFile = "opends.conf"
	defaultConfigType = "yaml"
	mainConfigPath    = "/etc/"
	altConfigPath     = "/etc/opends.d/"
	bindEnvOpts       = [][]string{
		{"public.host", "OPENDS_PUBLIC_HOST"},
		{"public.port", "OPENDS_PUBLIC_PORT"},
		{"admin.host", "OPENDS_ADMIN_HOST"},
		{"admin.port", "OPENDS_ADMIN_PORT"},
		{"database.driver", "OPENDS_DB_DRIVER"},
		{"database.dsn", "OPENDS_DB_DSN"},
		{"public.tls.enable", "OPENDS_PUBLIC_TLS_ENABLE"},
		{"public.tls.key_file", "OPENDS_PUBLIC_TLS_KEY_FILE"},
		{"public.tls.cert_file", "OPENDS_PUBLIC_TLS_CERT_FILE"},
		{"public.tls.key", "OPENDS_PUBLIC_TLS_KEY"},
		{"public.tls.cert", "OPENDS_PUBLIC_TLS_CERT"},
		{"logger.driver", "OPENDS_LOGGER_DRIVER"},
		{"logger.mode", "OPENDS_LOGGER_MODE"},
		{"cache.backend", "OPENDS_CACHE_BACKEND"},
		{"cache.address", "OPENDS_CACHE_ADDRESS"},
		{"cache.database", "OPENDS_CACHE_DB"},
		{"cors.enable", "OPENDS_CORS_ENABLE"},
		{"cors.max_age", "OPENDS_CORS_MAX_AGE"},
		{"cors.allowed_origins", "OPENDS_CORS_ALLOWED_ORIGINS"},
		{"cors.allowed_headers", "OPENDS_CORS_ALLOWED_HEADERS"},
		{"cors.allowed_methods", "OPENDS_CORS_ALLOWED_METHODS"},
		{"cors.allow_credentials", "OPENDS_CORS_ALLOW_CREDENTIALS"},
		{"debug", "OPENDS_DEBUG"},
	}
)

func NewConfig(path string) (*Config, error) {
	if path != "" {
		viper.SetConfigType(defaultConfigType)
		viper.SetConfigFile(path)
	} else {
		viper.SetConfigName(defaultConfigFile)
		viper.SetConfigType(defaultConfigType)
		viper.AddConfigPath(mainConfigPath)
		viper.AddConfigPath(altConfigPath)
		viper.AddConfigPath(".")
	}

	for _, env := range bindEnvOpts {
		viper.BindEnv(env...)
	}

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type Config struct {
	Public   ServerOptions   `mapstructure:"public"`
	Admin    ServerOptions   `mapstructure:"admin"`
	Database DatabaseOptions `mapstructure:"database"`
	Logger   LoggerOptions   `mapstructure:"logger"`
	Cache    CacheOptions    `mapstructure:"cache"`
	CORS     CORSOptions     `mapstructure:"cors"`
	Debug    bool            `mapstructure:"debug"`
}

type ServerOptions struct {
	Host string     `mapstructure:"host"`
	Port int16      `mapstructure:"port"`
	TLS  TLSOptions `mapstructure:"tls"`
}

type TLSOptions struct {
	Enable   bool   `mapstructure:"enable"`
	KeyFile  string `mapstructure:"key_file"`
	CertFile string `mapstructure:"cert_file"`
	Key      string `mapstructure:"key"`
	Cert     string `mapstructure:"cert"`
}

type DatabaseOptions struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"port"`
}

type LoggerOptions struct {
	Driver string `mapstructure:"driver"`
	Mode   string `mapstructure:"mode"`
}

type CacheOptions struct {
	Backend  string `mapstructure:"backend"`
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type CORSOptions struct {
	Enable           bool     `mapstructure:"enable"`
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	ExposedHeaders   []string `mapstructure:"exposed_headers"`
	MaxAge           int      `mapstructure:"max_age"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}
