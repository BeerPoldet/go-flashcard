package main

import "github.com/ilyakaznacheev/cleanenv"

type config struct {
	Server struct {
		Address string `toml:"address" env:"SERVER_ADDRESS" env-default:":3000"`
	} `toml:"server"`

	DB struct {
		ConnectionURI string `toml:"connectionURI" env:"DB_CONNECTION_URI"`
	} `toml:"db"`
}

func loadConfig() (*config, error) {
	var cfg config
	if err := cleanenv.ReadConfig("config.toml", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
