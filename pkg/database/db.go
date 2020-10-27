package database

import (
	"context"
	"errors"
	"fmt"
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
	useDatabasePool       bool
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
	QueryRowAndBind(ctx context.Context, sql string, args []interface{}, dst ...interface{}) error

	MustConnect(ctx context.Context)
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

func NewDB(dsn string, maxConnectionLifetime, maxIdleConnection, maxOpenConnection int64, useDatabasePool bool) IManager {
	return &database{
		dsn:                   dsn,
		maxConnectionLifetime: maxConnectionLifetime,
		maxIdleConnection:     maxIdleConnection,
		maxOpenConnection:     maxOpenConnection,
		useDatabasePool:       useDatabasePool,
	}
}

func (d *database) getConnection() {

}
func (d *database) Prepare(ctx context.Context, name, sql string) (sd *StatementDescription, err error) {
	if d.useDatabasePool {
		return nil, errors.New("prepare statement can't be used on connection pool mode")
	}
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
	if d.useDatabasePool {
		px, err = d.pool.Begin(ctx)
	} else {
		px, err = d.conn.Begin(ctx)
	}
	if err != nil {
		return nil, err
	}
	tx = *(*Tx)(&px)
	return
}

func (d *database) BeginTx(ctx context.Context, txOptions TxOptions) (tx Tx, err error) {
	var px pgx.Tx
	if d.useDatabasePool {
		px, err = d.pool.BeginTx(ctx, pgx.TxOptions(txOptions))
	} else {
		px, err = d.conn.BeginTx(ctx, pgx.TxOptions(txOptions))
	}
	if err != nil {
		return nil, err
	}
	tx = *(*Tx)(&px)
	return
}

func (d *database) Exec(ctx context.Context, sql string, arguments ...interface{}) (ct CommandTag, err error) {
	var pct pgconn.CommandTag
	if d.useDatabasePool {
		pct, err = d.pool.Exec(ctx, sql, arguments...)
	} else {
		pct, err = d.conn.Exec(ctx, sql, arguments...)
	}
	if err != nil {
		return nil, err
	}
	ct = *(*CommandTag)(&pct)
	return
}

func (d *database) Query(ctx context.Context, sql string, args ...interface{}) (rows Rows, err error) {
	var px pgx.Rows
	if d.useDatabasePool {
		px, err = d.pool.Query(ctx, sql, args...)
	} else {
		px, err = d.conn.Query(ctx, sql, args...)
	}
	if err != nil {
		return nil, err
	}
	rows = *(*Rows)(&px)
	return
}

func (d *database) QueryRow(ctx context.Context, sql string, args ...interface{}) Row {
	var row pgx.Row
	if d.useDatabasePool {
		row = d.pool.QueryRow(ctx, sql, args...)
	} else {
		row = d.conn.QueryRow(ctx, sql, args...)
	}
	return *&row
}

func (d *database) QueryRowAndBind(ctx context.Context, sql string, args []interface{}, dst ...interface{}) error {
	var err error
	if d.useDatabasePool {
		err = d.pool.QueryRow(ctx, sql, args...).Scan(dst...)
	} else {
		err = d.conn.QueryRow(ctx, sql, args...).Scan(dst...)
	}
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return fmt.Errorf("error: %s, %s", pgErr.Code, pgErr.Message)
		}
	}
	return nil
}

func (d *database) MustConnect(ctx context.Context) {
	if err := d.Connect(ctx); err != nil {
		log.Fatal(err)
	}
}

func (d *database) Connect(ctx context.Context) (err error) {
	if d.useDatabasePool {
		d.pool, err = pgxpool.Connect(ctx, d.dsn)
	} else {
		d.conn, err = pgx.Connect(ctx, d.dsn)
	}
	if err != nil {
		return err
	}
	return nil
}

func (d *database) Close(ctx context.Context) error {
	if d.useDatabasePool {
		d.pool.Close()
		return nil
	}
	return d.conn.Close(ctx)
}