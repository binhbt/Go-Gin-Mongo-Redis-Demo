package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis"
)

var (
	client = &redisClient{}
)

type redisClient struct {
	c *redis.Client
}

var redisClientInstance *redisClient

//GetClient get the redis client
func initialize() *redisClient {
	redisURI := os.Getenv("REDIS_URI")
	c := redis.NewClient(&redis.Options{
		Addr: redisURI,
	})

	if err := c.Ping().Err(); err != nil {
		panic("Unable to connect to redis " + err.Error())
	}
	client.c = c
	return client
}

//Get Redis Singleton Client
func getRedisClient() *redisClient {

	if redisClientInstance == nil {
		// <--- NOT THREAD SAFE /
		//Need optimize
		redisClientInstance = initialize()
	}

	return redisClientInstance
}

//GetKey get key
func (client *redisClient) getKey(key string, src interface{}) error {
	val, err := client.c.Get(key).Result()
	if err == redis.Nil || err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
		return err
	}
	return nil
}

//SetKey set key
func (client *redisClient) setKey(key string, value interface{}, expiration time.Duration) error {
	cacheEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = client.c.Set(key, cacheEntry, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
