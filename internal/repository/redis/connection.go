package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/env"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/repository"
	"github.com/kehl-gopher/movie-seat-reservation-theatre/internal/utility"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(config *env.Config) (*redis.Client, error) {

	p, err := utility.VerifyPort(config.Redis.RED_PORT)

	if err != nil {
		return nil, err
	}

	d, _ := strconv.Atoi(config.Redis.RED_DB)
	addr := fmt.Sprintf("%s:%d", config.Redis.RED_HOST, p)
	// Initialize Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       d,
		Password: config.Redis.RED_PASSWORD,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Println("failed to connect to redis database")
		return nil, err
	}

	fmt.Println("Redis connection successful")
	// Set the client to the repository
	repository.DB.Red = repository.NewRedis(client)
	return client, nil
}
