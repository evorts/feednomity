package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type redisManager struct {
	instance *redis.Client
	address  string
	password string
	db       int
}

func NewRedisStorage(address, password string, db int) IManager {
	return &redisManager{
		address:  address,
		password: password,
		db:       db,
	}
}

func (r *redisManager) Get(ctx context.Context, key string) interface{} {
	panic("implement me")
}

func (r *redisManager) Save(ctx context.Context, key string, value interface{}, expired int) error {
	panic("implement me")
}

func (r *redisManager) Delete(ctx context.Context, key string) error {
	panic("implement me")
}

func (r *redisManager) GetHash(ctx context.Context, key, hash string) interface{} {
	panic("implement me")
}

func (r *redisManager) GetHashTree(ctx context.Context, key string) map[string]interface{} {
	panic("implement me")
}

func (r *redisManager) SaveHash(ctx context.Context, key, hash string, value interface{}) error {
	panic("implement me")
}

func (r *redisManager) DeleteHashTree(ctx context.Context, key string) error {
	panic("implement me")
}

func (r *redisManager) DeleteHash(ctx context.Context, key, hash string) error {
	panic("implement me")
}

func (r *redisManager) MustConnect(ctx context.Context) {
	if err := r.Connect(ctx); err != nil {
		panic(err)
	}
}

func (r *redisManager) Connect(ctx context.Context) error {
	r.instance = redis.NewClient(&redis.Options{
		Addr:               r.address,
		Password:           r.password,
		DB:                 r.db,
	})
	cmd := r.instance.Ping(ctx)
	_, err := cmd.Result()
	return err
}
