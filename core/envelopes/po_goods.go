package envelopes

import (
	"github.com/shopspring/decimal"
	"imooc.com/resk/services"
	"time"
)

type RedEnvelopeGoods struct {
	Id             int64                `db:"id,omitempty"`
	EnvelopeNo     string               `db:"envelope_no,unique"`
	EnvelopeType   int                  `db:"envelope_type"`
	Username       string               `db:"username"`
	UserId         string               `db:"user_id"`
	Blessing       string               `db:"blessing"`
	Amount         decimal.Decimal      `db:"amount"`
	AmountOne      decimal.Decimal      `db:"amount_one"`
	Quantity       int                  `db:"quantity"`
	RemainAmount   decimal.Decimal      `db:"remain_amount"`
	RemainQuantity int                  `db:"remain_quantity"`
	ExpiredAt      time.Time            `db:"expired_at"`
	Status         services.OrderStatus `db:"status"`
	OrderType      services.OrderType   `db:"order_type"`
	PayStatus      services.PayStatus   `db:"pay_status"`
	CreatedAt      time.Time            `db:"created_at,omitempty"`
	UpdatedAt      time.Time            `db:"updated_at,omitempty"`
}

func (po *RedEnvelopeGoods) ToDTO() *services.RedEnvelopeGoodsDTO {
	return &services.RedEnvelopeGoodsDTO{
		EnvelopeNo:     po.EnvelopeNo,
		EnvelopeType:   po.EnvelopeType,
		Username:       po.Username,
		UserId:         po.UserId,
		Blessing:       po.Blessing,
		Amount:         po.Amount,
		AmountOne:      po.AmountOne,
		Quantity:       po.Quantity,
		RemainAmount:   po.RemainAmount,
		RemainQuantity: po.RemainQuantity,
		ExpiredAt:      po.ExpiredAt,
		Status:         int(po.Status),
		OrderType:      po.OrderType,
		PayStatus:      po.PayStatus,
		CreatedAt:      po.CreatedAt,
		UpdatedAt:      po.UpdatedAt,
		AccountNo:      "",
	}
}
