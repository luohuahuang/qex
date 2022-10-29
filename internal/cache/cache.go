package cache

import (
	"errors"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"strings"
	"time"
)

type Redis struct {
	Client *redis.Client
}

func New(endpoint string) *Redis {
	host := strings.Split(endpoint, "/")[0]
	db := strings.Split(endpoint, "/")[1]
	log.Printf("connect to redis: %s/%s", host, db)
	index, _ := strconv.Atoi(db)
	redisClient := &Redis{redis.NewClient(&redis.Options{
		Addr:         host,
		Password:     "",
		DB:           index,
		ReadTimeout:  10 * time.Second,
		DialTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})}
	return redisClient
}

func (r *Redis) Get(key string) (string, error) {
	log.Printf("get key: %s", key)
	results := r.Client.Get(key)
	if results != nil {
		value, err := r.Client.Get(key).Result()
		if err != nil {
			log.Println(err.Error()) // Redis `GET key` command. It returns redis.Nil error when key does not exist
			return "", nil
		}
		log.Printf("value: %s", value)
		return value, nil
	} else {
		return "", errors.New("fail to get cache")
	}
}

func (r *Redis) Set(key string, value interface{}, ttl int) error {
	log.Printf("Set Key:%s, Value:%v", key, value)
	err := r.Client.Set(key, value, time.Duration(ttl)*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) HSet(key string, field string, value string) error {
	log.Printf("HSet Key# %s, Value:%s", key, value)
	err := r.Client.HSet(key, field, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) HGet(key string, field string) (string, error) {
	log.Printf("HGet Key# %s, Field# %s", key, field)
	val, err := r.Client.HGet(key, field).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *Redis) Close() error {
	err := r.Client.Close()
	if err != nil {
		return err
	}
	return nil
}
