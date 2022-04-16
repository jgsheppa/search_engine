package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"log"
)

func NewServices() (*Services, error) {
	var pool = NewPool()

	client := pool.Get()
	defer client.Close()

	c, autocomplete, err := CreateIndex(pool, "articles")
	if err != nil {
		log.Println(err)
	}

	return &Services{
		Pool:         *pool,
		Redisearch:   *c,
		Autocomplete: *autocomplete,
	}, nil
}

type Services struct {
	Pool redis.Pool
	// Redisearch client
	Redisearch redisearch.Client
	// Redis autocomplete
	Autocomplete redisearch.Autocompleter
}

// NewPool creates the connection to the Redis DB
func NewPool() *redis.Pool {
	address, ok := viper.Get("REDIS_ADDRESS").(string)
	if !ok {
		log.Fatalf("Invalid type assertion for address")
	}
	password, ok := viper.Get("REDIS_PW").(string)
	if !ok {
		log.Fatalf("Invalid type assertion for password")
	}
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address,

				redis.DialPassword(password))
			if err != nil {
				log.Println(err)
			}
			return c, err
		},
	}
}
