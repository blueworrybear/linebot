package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Server Server
	Database Database
}

type Server struct {
	ChannelAccessToken string `envconfig:"CHANNEL_ACCESS_TOKEN"`
	ChannelSecret string `envconfig:"CHANNEL_SECRET"`
}

type Database struct {
	Driver      string `default:"sqlite3" envconfig:"DB_DRIVER"`
	Name        string `default:"core.db" envconfig:"DB_NAME"`
	User        string `envconfig:"DB_USER"`
	Password    string `envconfig:"DB_PASSWORD"`
	Socket   string `default:"/cloudsql" envconfig:"DB_SOCKET_DIR"`
	Instance string `envconfig:"INSTANCE_CONNECTION_NAME"`
	AutoMigrate bool   `default:"true" envconfig:"DB_AUTO_MIGRATE"`
}

func Environ() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	return cfg, err
}
