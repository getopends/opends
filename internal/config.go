package internal

type Config struct {
	Host string `mapstructure:"host"`
	Port int16  `mapstructure:"port"`
}
