package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
)

func signSecret(token string) (string, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()
	ctx := context.Background()
	signed := token + ":" + uuid.New().String()
	err := client.Set(ctx, token, signed, 36*10^11).Err()
	if err != nil {
		return "", err
	}
	return signed, nil
}

func updateSecret(token string, signed string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer client.Close()
	ctx := context.Background()
	err := client.Set(ctx, token, signed, 36*10^11).Err()
	if err != nil {
		return err
	}
	return nil
}

func checkSecret(token string) (bool, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
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
	return signed == parts[1], nil
}
