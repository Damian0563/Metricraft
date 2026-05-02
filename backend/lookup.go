package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
	"time"
)

func signToken(token string) (string, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis"),
		Password: "",
		DB:       0,
	})
	defer client.Close()
	ctx := context.Background()
	signed := token + ":" + uuid.New().String()
	err := client.Set(ctx, token, signed, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}
	return signed, nil
}

func updateToken(token string, signed string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis"),
		Password: "",
		DB:       0,
	})
	defer client.Close()
	ctx := context.Background()
	err := client.Set(ctx, token, signed, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func checkToken(token string) (bool, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis"),
		Password: "",
		DB:       0,
	})
	defer client.Close()
	ctx := context.Background()
	parts := strings.Split(token, ":")
	signed, err := client.Get(ctx, parts[0]).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	if len(parts) == 1 {
		return false, nil
	}
	return signed == token, nil
}
