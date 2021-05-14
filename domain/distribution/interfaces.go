package distribution

import "context"

type IManager interface {
	FindByIds(ctx context.Context, ids ...int64) ([]*Distribution, error)
	FindAll(ctx context.Context, page, limit int) (items []*Distribution, total int, err error)
	InsertMultiple(ctx context.Context, items []*Distribution) error
	Update(ctx context.Context, item Distribution) error
	UpdateStatusAndCountByIds(ctx context.Context, ids ...int64) error
	DeleteByIds(ctx context.Context, ids ...int64) error

	FindObjectByIds(ctx context.Context, ids ...int64) ([]*Object, error)
	FindObjectByLinkIds(ctx context.Context, ids ...int64) ([]*Object, error)
	FindObjectByRespondentAndLinkId(ctx context.Context, respondentId, id int64) ([]*Object, error)
	FindObjectsByDistributionIds(ctx context.Context, ids ...int64) ([]*Object, error)
	FindAllObjects(ctx context.Context, page, limit int) (items []*Object, total int, err error)
	InsertObjects(ctx context.Context, items []*Object) error
	UpdateObject(ctx context.Context, item Object) error
	UpdateObjectStatusByIds(ctx context.Context, status PublishingStatus, ids ...int64) error
	UpdateObjectRetryCountByIds(ctx context.Context, ids ...int64) error
	DeleteObjectByIds(ctx context.Context, ids ...int64) error

	FindAllQueues(ctx context.Context, page, limit int) (items []*Queue, total int64, err error)
	InsertQueues(ctx context.Context, items []*Queue) ([]int64, error)
	DeleteQueueByIds(ctx context.Context, ids ...int64) error

	InsertLogs(ctx context.Context, items []*Log) error
}
