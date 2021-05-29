package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
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

func (r *redisManager) Get(ctx context.Context, key string) (interface{}, error) {
	cmd := r.instance.Get(ctx, key)
	err := cmd.Err()
	if err != nil {
		return nil, err
	}
	var rs interface{}
	err = cmd.Scan(&rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (r *redisManager) GetString(ctx context.Context, key string, def string) string {
	cmd := r.instance.Get(ctx, key)
	if cmd.Err() != nil {
		return def
	}
	return cmd.String()
}

func (r *redisManager) GetInt64(ctx context.Context, key string, def int64) int64 {
	rs, err := r.instance.Get(ctx, key).Int64()
	if err != nil {
		return def
	}
	return rs
}

func (r *redisManager) GetInt(ctx context.Context, key string, def int) int {
	rs, err := r.instance.Get(ctx, key).Int()
	if err != nil {
		return def
	}
	return rs
}

func (r *redisManager) GetTime(ctx context.Context, key string, def *time.Time) *time.Time {
	rs, err := r.instance.Get(ctx, key).Time()
	if err != nil {
		return def
	}
	return &rs
}

func (r *redisManager) GetBool(ctx context.Context, key string, def bool) bool {
	rs, err := r.instance.Get(ctx, key).Bool()
	if err != nil {
		return def
	}
	return rs
}

func (r *redisManager) GetFloat64(ctx context.Context, key string, def float64) float64 {
	rs, err := r.instance.Get(ctx, key).Float64()
	if err != nil {
		return def
	}
	return rs
}

func (r *redisManager) GetFloat32(ctx context.Context, key string, def float32) float32 {
	rs, err := r.instance.Get(ctx, key).Float32()
	if err != nil {
		return def
	}
	return rs
}

func (r *redisManager) Set(ctx context.Context, key string, value interface{}, expired int64) error {
	cmd := r.instance.Set(ctx, key, value, time.Duration(expired) * time.Second)
	return cmd.Err()
}

func (r *redisManager) Delete(ctx context.Context, key string) error {
	cmd := r.instance.Del(ctx, key)
	return cmd.Err()
}

func (r *redisManager) GetHash(ctx context.Context, key, hash string) (interface{}, error) {
	cmd := r.instance.HGet(ctx, key, hash)
	err := cmd.Err()
	if err != nil {
		return nil, err
	}
	var rs interface{}
	err = cmd.Scan(&rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (r *redisManager) GetHashTree(ctx context.Context, key string) (map[string]interface{}, error) {
	cmd := r.instance.HGetAll(ctx, key)
	err := cmd.Err()
	if err != nil {
		return nil, err
	}
	var rs map[string]interface{}
	err = cmd.Scan(&rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (r *redisManager) SetHash(ctx context.Context, key string, values ...interface{}) error {
	cmd := r.instance.HSet(ctx, key, values)
	return cmd.Err()
}

func (r *redisManager) DeleteHashTree(ctx context.Context, key string) error {
	cmd := r.instance.Del(ctx, key)
	return cmd.Err()
}

func (r *redisManager) DeleteHash(ctx context.Context, key, hash string) error {
	cmd := r.instance.HDel(ctx, key, hash)
	return cmd.Err()
}

func (r *redisManager) MustConnect(ctx context.Context) {
	if err := r.Connect(ctx); err != nil {
		log.Fatal(err)
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
