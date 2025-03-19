package redis

import (
	"fmt"

	"github.com/labstack/gommon/log"
	redislib "github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}

type Adapter struct {
	client *redislib.Client
}

func New(config Config) *Adapter {
	rdb := redislib.NewClient(&redislib.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
	log.Info("✅ Redis is up running...")

	return &Adapter{client: rdb}
}

func (a Adapter) Client() *redislib.Client {
	return a.client
}
