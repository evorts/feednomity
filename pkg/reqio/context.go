package reqio

import (
	"context"
)

type Context struct {
	ctx context.Context
}

type IContext interface {
	Get(key string) interface{}
	Value() context.Context
}

func NewContext(ctx context.Context) IContext {
	return &Context{ctx: ctx}
}

func (c *Context) Value() context.Context {
	return c.ctx
}

func (c *Context) Get(key string) interface{} {
	return c.ctx.Value(key)
}
