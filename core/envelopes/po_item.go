package envelopes

import (
	"github.com/shopspring/decimal"
	"imooc.com/resk/services"
	"time"
)

type RedEnvelopeItem struct {
	Id           int64           `db:"id,omitempty"`
	ItemNo       int64           `db:"item_no,unique"`       // 红包订单详情编号
	EnvelopeNo   string          `db:"envelope_no"`          // 红包编号
	RecvUsername string          `db:"recv_username"`        // 接收者用户名
	RecvUserId   string          `db:"recv_user_id"`         // 接收者用户id
	Amount       decimal.Decimal `db:"amount"`               // 收到金额
	Quantity     int             `db:"quantity"`             // 收到数量
	RemainAmount decimal.Decimal `db:"remain_amount"`        // 剩余金额
	AccountNo    string          `db:"account_no"`           // 红包接收者账户ID
	PayStatus    int             `db:"pay_status"`           // 支付状态
	CreatedAt    time.Time       `db:"created_at,omitempty"` // 创建时间
	UpdatedAt    time.Time       `db:"updated_at,omitempty"` // 修改时间
}

func (po *RedEnvelopeItem) ToDTO() *services.RedEnvelopeItemDTO {
	return &services.RedEnvelopeItemDTO{
		ItemNo:       po.ItemNo,
		EnvelopeNo:   po.EnvelopeNo,
		RecvUsername: po.RecvUsername,
		RecvUserId:   po.RecvUserId,
		Amount:       po.Amount,
		Quantity:     po.Quantity,
		RemainAmount: po.RemainAmount,
		AccountNo:    po.AccountNo,
		PayStatus:    po.PayStatus,
		CreatedAt:    po.CreatedAt,
		UpdatedAt:    po.UpdatedAt,
	}
}
