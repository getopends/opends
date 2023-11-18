package internal

import (
	"fmt"

	"github.com/spf13/viper"
)

/*

{
   "public": {
   		"host": "127.0.0.1",
        "port": 13000
   },
   "admin": {
   		"host": "127.0.0.1",
        "port": 13001
   },
   "logger": {
   		"driver": "zap",
        "mode": "prod"
   },
   "database": {
   		"driver": "postgres",
        "dsn": "postgres@postgres:5432/db"
   }
}


*/

var (
	envPrefix   = "OPENDS"
	bindEnvOpts = []BindEnvOptions{
		{
			Key:   "Public.Host",
			Value: "PUBLIC_HOST",
		},
		{
			Key:   "Public.Port",
			Value: "PUBLIC_PORT",
		},
		{
			Key:   "Admin.Host",
			Value: "ADMIN_HOST",
		},
		{
			Key:   "Admin.Port",
			Value: "ADMIN_PORT",
		},
		{
			Key:   "Database.Driver",
			Value: "DB_DRIVER",
		},
		{
			Key:   "Database.DSN",
			Value: "DB_DSN",
		},
		{
			Key:   "Public.TLS.Enable",
			Value: "PUBLIC_TLS_ENABLE",
		},
		{
			Key:   "Public.TLS.KeyFile",
			Value: "PUBLIC_TLS_KEY_FILE",
		},
		{
			Key:   "Public.TLS.CertFile",
			Value: "PUBLIC_TLS_CERT_FILE",
		},
		{
			Key:   "Public.TLS.Key",
			Value: "PUBLIC_TLS_KEY",
		},
		{
			Key:   "Public.TLS.Cert",
			Value: "PUBLIC_TLS_CERT",
		},
		{
			Key:   "Logger.Driver",
			Value: "LOGGER_DRIVER",
		},
		{
			Key:   "Logger.Mode",
			Value: "LOGGER_MODE",
		},
		{
			Key:   "Cache.Backend",
			Value: "CACHE_BACKEND",
		},
		{
			Key:   "Cache.Address",
			Value: "CACHE_ADDRESS",
		},
		{
			Key:   "Cache.Database",
			Value: "CACHE_DB",
		},
		{
			Key:   "CORS.Enable",
			Value: "CORS_ENABLE",
		},
		{
			Key:   "CORS.MaxAge",
			Value: "CORS_MAX_AGE",
		},
		{
			Key:   "CORS.AllowedOrigns",
			Value: "CORS_ALLOWED_ORIGINS",
		},
		{
			Key:   "CORS.AllowedHeaders",
			Value: "CORS_ALLOWED_HEADERS",
		},
		{
			Key:   "CORS.AllowedMethods",
			Value: "CORS_ALLOWED_METHODS",
		},
		{
			Key:   "CORS.ExposedHeaders",
			Value: "CORS_EXPOSED_METHODS",
		},
		{
			Key:   "CORS.AllowCredentials",
			Value: "CORS_ALLOW_CREDENTIALS",
		},
	}
)

func ConfigInit() (*Config, error) {
	viper.SetConfigName("opends.conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/opends.conf")
	viper.AddConfigPath("/etc/opends.d/opends.conf")
	viper.AddConfigPath(".")

	for _, env := range bindEnvOpts {
		viper.BindEnv(env.Key, fmt.Sprintf("%v_%v", envPrefix, env.Value))
	}

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	cfg := &Config{}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

type Config struct {
	Public   ServerOptions   `mapstructure:"public"`
	Admin    ServerOptions   `mapstructure:"admin"`
	Database DatabaseOptions `mapstructure:"database"`
	Logger   LoggerOptions   `mapstructure:"logger"`
	Cache    CacheOptions    `mapstructure:"cache"`
	CORS     CORSOptions     `mapstructure:"cors"`
}

type ServerOptions struct {
	Host string     `mapstructure:"host"`
	Port int16      `mapstructure:"port"`
	TLS  TLSOptions `mapstructure:"tls"`
}

type TLSOptions struct {
	Enable   bool   `mapstructure:"enable"`
	KeyFile  string `mapstructure:"key-file"`
	CertFile string `mapstructure:"cert-file"`
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
	AllowedOrigins   []string `mapstructure:"allowed-origins"`
	AllowedMethods   []string `mapstructure:"allowed-methods"`
	AllowedHeaders   []string `mapstructure:"allowed-headers"`
	ExposedHeaders   []string `mapstructure:"exposed-headers"`
	MaxAge           int      `mapstructure:"max-age"`
	AllowCredentials bool     `mapstructure:"allow-credentials"`
}

type BindEnvOptions struct {
	Key   string
	Value string
}
