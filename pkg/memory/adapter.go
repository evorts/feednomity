package memory

import "context"

type IManager interface {
	Get(ctx context.Context, key string) interface{}
	Save(ctx context.Context, key string, value interface{}, expired int) error
	Delete(ctx context.Context, key string) error

	GetHash(ctx context.Context, key, hash string) interface{}
	GetHashTree(ctx context.Context, key string) map[string]interface{}
	SaveHash(ctx context.Context, key, hash string, value interface{}) error
	DeleteHashTree(ctx context.Context, key string) error
	DeleteHash(ctx context.Context, key, hash string) error

	MustConnect(ctx context.Context)
	Connect(ctx context.Context) error
}
