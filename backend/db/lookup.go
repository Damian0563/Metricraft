package db

import (
	"backend/types"
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func SignToken(token string, rotate bool) (string, error) {
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
	err := client.Set(ctx, token, signed, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}
	return signed, nil
}

func UpdateToken(token string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis"),
		Password: "",
		DB:       0,
	})
	defer client.Close()
	ctx := context.Background()
	parts := strings.Split(token, ":")
	err := client.Set(ctx, parts[0], token, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func CheckToken(token string) (bool, error) {
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

func GenerateCode() string {
	var code string
	var index int
	for index < 6 {
		code += strconv.Itoa(rand.Intn(10))
		index++
	}
	return code
}

func SetCodeValidity(mail string, code string) error {
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

func CheckCodeValidity(routine chan types.ExistsErrResponse, mail string, code string) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis"),
		Password: "",
		DB:       0,
	})
	var origin string = "checkCodeValidity"
	defer client.Close()
	ctx := context.Background()
	key := "verify:" + mail
	signed, err := client.Get(ctx, key).Result()
	if err != nil || signed == "" {
		if err == redis.Nil {
			routine <- types.ExistsErrResponse{Exists: false, Err: nil, Origin: origin}
			return
		}
		routine <- types.ExistsErrResponse{Exists: false, Err: err, Origin: origin}
		return
	}
	routine <- types.ExistsErrResponse{Exists: signed == code, Err: nil, Origin: origin}
}
