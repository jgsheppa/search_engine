package models

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewServices(connectionInfo string) (*Services, error) {
	var pool = NewPool()

	client := pool.Get()
	defer client.Close()

	c, autocomplete, err := CreateIndex(pool, "index")
	if err != nil {
		log.Println(err)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}

	return &Services{
		Pool:         *pool,
		Redisearch:   *c,
		Autocomplete: *autocomplete,
		User:         NewUserService(db),
		db:           db,
	}, nil
}

type Services struct {
	Pool redis.Pool
	// Redisearch client
	Redisearch redisearch.Client
	// Redis autocomplete
	Autocomplete redisearch.Autocompleter
	User         UserService
	db           *gorm.DB
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

// Used to delete old database tables and entries in development
func (s *Services) DestructiveReset() error {
	err := s.db.Migrator().DropTable(&User{})
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automigrate all database tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{})
}
