package redis

import (
	"context"
	"log"
	"time"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	goredis "github.com/redis/go-redis/v9"
)

var Client *goredis.Client

func Conn() {
	Client = goredis.NewClient(&goredis.Options{
		Addr:     dotenv.RedisHost + ":" + dotenv.RedisPort,
		Username: "default",
		Password: dotenv.RedisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		log.Panicf("error connecting to redis: %v", err)
	}

	log.Println("Redis connected successfully")
}
