package memory

import (
	"context"
	"time"
)

type IManager interface {
	Get(ctx context.Context, key string) (interface{}, error)
	GetString(ctx context.Context, key string, def string) string
	GetInt64(ctx context.Context, key string, def int64) int64
	GetInt(ctx context.Context, key string, def int) int
	GetTime(ctx context.Context, key string, def *time.Time) *time.Time
	GetBool(ctx context.Context, key string, def bool) bool
	GetFloat64(ctx context.Context, key string, def float64) float64
	GetFloat32(ctx context.Context, key string, def float32) float32

	Set(ctx context.Context, key string, value interface{}, expired int64) error
	Delete(ctx context.Context, key string) error

	GetHash(ctx context.Context, key, hash string) (interface{}, error)
	GetHashTree(ctx context.Context, key string) (map[string]interface{}, error)
	SetHash(ctx context.Context, key string, values ...interface{}) error
	DeleteHashTree(ctx context.Context, key string) error
	DeleteHash(ctx context.Context, key, hash string) error

	MustConnect(ctx context.Context)
	Connect(ctx context.Context) error
}
