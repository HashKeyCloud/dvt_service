package store

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"DVT_Service/conf"
)

// Store Service Persistence Related
type Store struct {
	mysql *gorm.DB
	redis *redis.Client
}

// NewStore New Store struct
func NewStore(cfg *conf.ConfigDB) *Store {
	mysql, err := mysqlConnect(cfg.Mysql)
	if err != nil {
		panic(err)
	}

	rds := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := rds.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return &Store{mysql: mysql, redis: rds}
}
