package internal

import (
	"fmt"
	"log"

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
			Key:   "public.host",
			Value: "PUBLIC_HOST",
		},
		{
			Key:   "public.port",
			Value: "PUBLIC_PORT",
		},
		{
			Key:   "admin.host",
			Value: "ADMIN_HOST",
		},
		{
			Key:   "admin.port",
			Value: "ADMIN_PORT",
		},
		{
			Key:   "database.driver",
			Value: "DB_DRIVER",
		},
		{
			Key:   "database.dsn",
			Value: "DB_DSN",
		},
		{
			Key:   "public.tls.enable",
			Value: "PUBLIC_TLS_ENABLE",
		},
		{
			Key:   "public.tls.key_file",
			Value: "PUBLIC_TLS_KEY_FILE",
		},
		{
			Key:   "public.tls.cert_file",
			Value: "PUBLIC_TLS_CERT_FILE",
		},
		{
			Key:   "public.tls.key",
			Value: "PUBLIC_TLS_KEY",
		},
		{
			Key:   "public.tls.cert",
			Value: "PUBLIC_TLS_CERT",
		},
		{
			Key:   "logger.driver",
			Value: "LOGGER_DRIVER",
		},
		{
			Key:   "logger.mode",
			Value: "LOGGER_MODE",
		},
		{
			Key:   "cache.backend",
			Value: "CACHE_BACKEND",
		},
		{
			Key:   "cache.address",
			Value: "CACHE_ADDRESS",
		},
		{
			Key:   "cache.database",
			Value: "CACHE_DB",
		},
		{
			Key:   "cors.enable",
			Value: "CORS_ENABLE",
		},
		{
			Key:   "cors.max_age",
			Value: "CORS_MAX_AGE",
		},
		{
			Key:   "cors.allowed_origins",
			Value: "CORS_ALLOWED_ORIGINS",
		},
		{
			Key:   "cors.allowed_headers",
			Value: "CORS_ALLOWED_HEADERS",
		},
		{
			Key:   "cors.allowed_methods",
			Value: "CORS_ALLOWED_METHODS",
		},
		{
			Key:   "cors.allow_credentials",
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
		log.Printf("Binding %v to %v\n", env.Value, env.Key)

		viper.BindEnv(env.Key, fmt.Sprintf("%v_%v", envPrefix, env.Value))
	}

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	log.Println(viper.AllKeys())

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
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	ExposedHeaders   []string `mapstructure:"exposed_headers"`
	MaxAge           int      `mapstructure:"max_age"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

type BindEnvOptions struct {
	Key   string
	Value string
}
