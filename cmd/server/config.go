package main

import "github.com/ilyakaznacheev/cleanenv"

type config struct {
	Server struct {
		Address string `toml:"address" env:"SERVER_ADDRESS" env-default:":3000"`
	} `toml:"server"`
}

func loadConfig() (*config, error) {
	var cfg config
	if err := cleanenv.ReadConfig("config.toml", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
