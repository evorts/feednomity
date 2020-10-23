package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type database struct {
	dsn                   string
	maxConnectionLifetime int64
	maxIdleConnection     int64
	maxOpenConnection     int64
	conn                  *pgx.Conn
	pool                  *pgxpool.Pool
}

type StatementDescription pgconn.StatementDescription
type Tx pgx.Tx
type TxOptions pgx.TxOptions
type CommandTag pgconn.CommandTag
type Rows pgx.Rows
type Row pgx.Row

type IManager interface {
	Prepare(ctx context.Context, name, sql string) (sd *StatementDescription, err error)
	Begin(ctx context.Context) (Tx, error)
	BeginTx(ctx context.Context, txOptions TxOptions) (Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row

	MustConnect(ctx context.Context)
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

func NewDB(dsn string, maxConnectionLifetime, maxIdleConnection, maxOpenConnection int64) IManager {
	return &database{
		dsn:                   dsn,
		maxConnectionLifetime: maxConnectionLifetime,
		maxIdleConnection:     maxIdleConnection,
		maxOpenConnection:     maxOpenConnection,
	}
}

func (d *database) Prepare(ctx context.Context, name, sql string) (sd *StatementDescription, err error) {
	var rs *pgconn.StatementDescription
	rs, err = d.conn.Prepare(ctx, name, sql)
	if err != nil {
		return nil, err
	}
	sd = (*StatementDescription)(rs)
	return
}

func (d *database) Begin(ctx context.Context) (tx Tx, err error) {
	var px pgx.Tx
	px, err = d.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	tx = *(*Tx)(&px)
	return
}

func (d *database) BeginTx(ctx context.Context, txOptions TxOptions) (tx Tx, err error) {
	var px pgx.Tx
	px, err = d.conn.BeginTx(ctx, pgx.TxOptions(txOptions))
	if err != nil {
		return nil, err
	}
	tx = *(*Tx)(&px)
	return
}

func (d *database) Exec(ctx context.Context, sql string, arguments ...interface{}) (ct CommandTag, err error) {
	var pct pgconn.CommandTag
	pct, err = d.conn.Exec(ctx, sql, arguments...)
	if err != nil {
		return nil, err
	}
	ct = *(*CommandTag)(&pct)
	return
}

func (d *database) Query(ctx context.Context, sql string, args ...interface{}) (rows Rows, err error) {
	var px pgx.Rows
	px, err = d.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	rows = *(*Rows)(&px)
	return
}

func (d *database) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	row := d.conn.QueryRow(ctx, sql, args...)
	return *(*Row)(&row)
}

func (d *database) MustConnect(ctx context.Context) {
	if err := d.Connect(ctx); err != nil {
		log.Fatal(err)
	}
}

func (d *database) Connect(ctx context.Context) (err error) {
	d.pool, err = pgxpool.Connect(ctx, d.dsn)
	if err != nil {
		return err
	}
	if err = d.conn.Ping(ctx); err != nil {
		return err
	}
	return nil
}

func (d *database) Close(ctx context.Context) error {
	return d.conn.Close(ctx)
}
