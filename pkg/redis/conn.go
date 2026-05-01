package redis

import (
	"context"
	"log"
	"time"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	goredis "github.com/redis/go-redis/v9"
)

var Client *goredis.Client

func Conn(cfg dotenv.RedisConfig) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Username: "default",
		Password: cfg.RedisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("Redis connected successfully")
	return client, nil
}
