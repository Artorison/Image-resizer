package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger LoggerConf `yaml:"logger"`
	App    App        `yaml:"app"`
}

type HTTPServer struct {
	Port string `yaml:"port" env:"APP_PORT" env-default:"8080"`
}
type LoggerConf struct {
	Level string `yaml:"level"`
	Type  string `yaml:"type"`
}

type App struct {
	CacheSize int           `yaml:"cache_size" env:"CACHE_SIZE" env-default:"10"`
	Port      string        `yaml:"port" env:"APP_PORT" env-default:"8080"`
	Timeout   time.Duration `yaml:"timeout"`
	Quality   int           `yaml:"image_quality"`
}

func MustLoad(filepath string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(filepath, &cfg)
	if err != nil {
		panic("failed to read config from file: " + filepath)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to load envs")
	}

	return &cfg
}
