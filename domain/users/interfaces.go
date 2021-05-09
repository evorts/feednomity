package users

import "context"

type IUsers interface {
	FindByIds(ctx context.Context, ids ...int64) ([]*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByNameAndGroupId(ctx context.Context, name string, groupId int) ([]*User, error)
	FindByNameAndOrgId(ctx context.Context, name string, orgId int) ([]*User, error)
	FindByName(ctx context.Context, name string) ([]*User, error)
	FindAll(ctx context.Context, page, limit int) (items []*User, total int, err error)

	Insert(ctx context.Context, item User) error
	InsertMultiple(ctx context.Context, items []*User) error
	Update(ctx context.Context, item User) error
	DeleteByIds(ctx context.Context, id []int64) error
	DisableByIds(ctx context.Context, id []int64) error

	FindGroupByIds(ctx context.Context, ids ...int64) ([]*Group, error)
	FindGroupByOrgId(ctx context.Context, id int64) ([]*Group, error)
	FindAllGroups(ctx context.Context, page, limit int) (items []*Group, total int, err error)
	InsertGroup(ctx context.Context, item Group) error
	InsertGroups(ctx context.Context, items []*Group) error
	UpdateGroup(ctx context.Context, item Group) error
	DeleteGroupByIds(ctx context.Context, ids ...int64) error

	FindOrganizationByIds(ctx context.Context, ids ...int64) ([]*Organization, error)
	FindAllOrganizations(ctx context.Context, page, limit int) (items []*Organization, total int, err error)
	InsertOrganization(ctx context.Context, item Organization) error
	InsertOrganizations(ctx context.Context, items []*Organization) error
	UpdateOrganization(ctx context.Context, item Organization) error
	DeleteOrganizationByIds(ctx context.Context, ids ...int64) error
}
