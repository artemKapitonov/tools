package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Config is type for database connection
type Config struct {
	Host     string `yaml:"host" env:"REDIS_HOST"      env-default:"localhost"`
	Port     int    `yaml:"port" env:"REDIS_PORT"      env-default:"6379"`
	DBNumber int    `yaml:"db"   env:"REDIS_DB_NUMBER" env-default:"1"`
	User     string `yaml:"user" env:"REDIS_USER"      env-default:"admin"`
	Password string `env:"REDIS_USER_PASSWORD"`
}

func getConfig(cfg *Config) *redis.Options {
	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DBNumber,
	}
}

// ConnectToDb creates connection
func ConnectToRedis(ctx context.Context, cfg *Config) (Client, error) {
	config := getConfig(cfg)
	
	client := redis.NewClient(config)

	if err := client.Ping(ctx).Err(); err != nil {
		return Client{}, err
	}

	return Client(*client), nil
}
