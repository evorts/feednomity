package distribution

import "context"

type IManager interface {
	FindByIds(ctx context.Context, ids ...int64) ([]*Distribution, error)
	FindAll(ctx context.Context, page, limit int) (items []*Distribution, total int, err error)
	InsertMultiple(ctx context.Context, items []*Distribution) error
	Update(ctx context.Context, item Distribution) error
	DeleteByIds(ctx context.Context, ids ...int64) error

	FindObjectByIds(ctx context.Context, ids ...int64) ([]*Object, error)
	FindAllObjects(ctx context.Context, page, limit int) (items []*Object, total int, err error)
	InsertObjects(ctx context.Context, items []*Object) error
	UpdateObject(ctx context.Context, item Object) error
	DeleteObjectByIds(ctx context.Context, ids ...int64) error
}
