package test

import (
	"github.com/shopspring/decimal"
	"time"
)

// 红包商品表, 默认有符号类型字段
type GoodsSigned struct {
	Goods
}

// 红包商品表，无符号类型字段
type GoodsUnsigned struct {
	Goods
}

// 红包商品
type Goods struct {
	RemainAmount   decimal.Decimal `db:"remain_amount"`
	RemainQuantity int             `db:"remain_quantity"`
	EnvelopeNo     string          `db:"envelope_no,uni"`
	CreatedAt      time.Time       `db:"created_at,omitempty"`
	UpdatedAt      time.Time       `db:"updated_at,omitempty"`
}
