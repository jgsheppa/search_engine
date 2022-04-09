package redis_conn

import (
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"log"
)

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
