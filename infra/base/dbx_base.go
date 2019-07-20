package base

import (
	"context"
	"database/sql"
	"github.com/tietang/dbx"
)

const TX = "tx"

type BaseDao struct {
	Tx *sql.Tx
}

func (d *BaseDao) SetTx(tx *sql.Tx) {
	d.Tx = tx
}

type TxFunc func(runner *dbx.TxRunner) error

// 事务执行
func Tx(fn TxFunc) error {
	return TxContext(context.Background(), fn)
}

// 实物执行
func TxContext(ctx context.Context, fn TxFunc) error {
	return DbxDatabase().Tx(fn)
}

// 将runner绑定到上下文
func WithValueContext(parent context.Context, runner *dbx.TxRunner) context.Context {
	return context.WithValue(parent, TX, runner)
}

func ExcuteContext(ctx context.Context, fn func(runner *dbx.TxRunner) error) error {
	tx := ctx.Value(TX).(*dbx.TxRunner)
	return fn(tx)
}
