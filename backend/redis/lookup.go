package redis

import (
	"backend/types"
	"context"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	client     *goredis.Client
	clientOnce sync.Once
)

func InitClient() {
	clientOnce.Do(func() {
		client = goredis.NewClient(&goredis.Options{
			Addr:     os.Getenv("redis"),
			Password: "",
			DB:       0,
		})
	})
}

func Client() *goredis.Client {
	InitClient()
	return client
}

func SignToken(token string, rotate bool) (string, error) {
	ctx := context.Background()
	redisClient := Client()
	if !rotate {
		return redisClient.Get(ctx, token).Result()
	}
	signed := token + ":" + uuid.New().String()
	err := redisClient.Set(ctx, token, signed, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}
	return signed, nil
}

func SetRecovery(token string, recovery string) error {
	ctx := context.Background()
	redisClient := Client()
	err := redisClient.Set(ctx, "recovery:"+recovery, token, 10*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func FetchRecoveryUser(recovery string) (string, error) {
	ctx := context.Background()
	redisClient := Client()
	token, err := redisClient.Get(ctx, "recovery:"+recovery).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

func InvalidateRecovery(recovery string) error {
	ctx := context.Background()
	redisClient := Client()
	err := redisClient.Del(ctx, "recovery:"+recovery).Err()
	if err != nil {
		return err
	}
	return nil
}

func InvalidateSession(token string) error {
	ctx := context.Background()
	redisClient := Client()
	err := redisClient.Del(ctx, token).Err()
	if err != nil {
		return err
	}
	return nil
}

func UpdateToken(token string) error {
	ctx := context.Background()
	redisClient := Client()
	parts := strings.Split(token, ":")
	err := redisClient.Set(ctx, parts[0], token, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func CheckToken(token string) (bool, error) {
	ctx := context.Background()
	redisClient := Client()
	parts := strings.Split(token, ":")
	signed, err := redisClient.Get(ctx, parts[0]).Result()
	if err != nil {
		if err == goredis.Nil {
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
	ctx := context.Background()
	redisClient := Client()
	key := "verify:" + mail
	err := redisClient.Set(ctx, key, code, 11*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func CheckCodeValidity(routine chan types.ExistsErrResponse, mail string, code string) {
	var origin string = "checkCodeValidity"
	ctx := context.Background()
	redisClient := Client()
	key := "verify:" + mail
	signed, err := redisClient.Get(ctx, key).Result()
	if err != nil || signed == "" {
		if err == goredis.Nil {
			routine <- types.ExistsErrResponse{Exists: false, Err: nil, Origin: origin}
			return
		}
		routine <- types.ExistsErrResponse{Exists: false, Err: err, Origin: origin}
		return
	}
	routine <- types.ExistsErrResponse{Exists: signed == code, Err: nil, Origin: origin}
}
