package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"os"
	"strings"
	"time"
	"strconv"
)

func signToken(token string, rotate bool) (string, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis"),
		Password: "",
		DB:       0,
	})
	defer client.Close()
	ctx := context.Background()
	if !rotate {
		return client.Get(ctx, token).Result()
	}
	signed := token + ":" + uuid.New().String()
	err := client.Set(ctx, token, signed, 1*time.Hour).Err()
	if err != nil {
		return "", err
	}
	return signed, nil
}

func updateToken(token string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis"),
		Password: "",
		DB:       0,
	})
	defer client.Close()
	ctx := context.Background()
	parts := strings.Split(token, ":")
	err := client.Set(ctx, parts[0], token, 1*time.Hour).Err()
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

func generateCode() string {
	var code string
	var index int
	for index < 6 {
		code += strconv.Itoa(rand.Intn(10))
		index++
	}
	return code
}

func setCodeValidity(mail string, code string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis"),
		Password: "",
		DB:       0,
	})
	defer client.Close()
	ctx := context.Background()
	key := "verify:" + mail
	err := client.Set(ctx, key, code, 11*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func checkCodeValidity(mail string, code string) (bool, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis"),
		Password: "",
		DB:       0,
	})
	defer client.Close()
	ctx := context.Background()
	key := "verify:" + mail
	signed, err := client.Get(ctx, key).Result()
	if err != nil || signed == "" {
		return false, err
	}
	return signed == code, nil
}
