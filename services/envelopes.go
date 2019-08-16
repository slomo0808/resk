package services

import (
	"github.com/shopspring/decimal"
	"imooc.com/resk/infra/base"
	"time"
)

var IRedEnvelopeService RedEnvelopeService

// 用于对外暴露的应用服务，唯一的暴露点
func GetRedEnvelopeService() RedEnvelopeService {
	base.Check(IRedEnvelopeService)
	return IRedEnvelopeService
}

// 创建接口
type RedEnvelopeService interface {
	// 发红包
	SendOut(RedEnvelopeSendingDTO) (*RedEnvelopeActivity, error)
	// 收红包
	Receive(*RedEnvelopeReceiveDTO) (*RedEnvelopeItemDTO, error)
	// 退款
	Refund(string) *RedEnvelopeGoodsDTO
	// 查询红包订单
	Get() RedEnvelopeGoodsDTO
}

// 发红包所需信息
type RedEnvelopeSendingDTO struct {
	EnvelopeType int             `json:"envelope_type" validate:"required"`
	Username     string          `json:"username" validate:"required"`
	UserId       string          `json:"user_id" validate:"required"`
	Blessing     string          `json:"blessing"`
	Amount       decimal.Decimal `json:"amount" validate:"required"`
	Quantity     int             `json:"quantity" validate:"required"`
}

func (dto *RedEnvelopeSendingDTO) ToGoods() *RedEnvelopeGoodsDTO {
	return &RedEnvelopeGoodsDTO{
		EnvelopeType: dto.EnvelopeType,
		Username:     dto.Username,
		UserId:       dto.UserId,
		Blessing:     dto.Blessing,
		Amount:       dto.Amount,
		Quantity:     dto.Quantity,
	}
}

// 收红包所需信息
type RedEnvelopeReceiveDTO struct {
	EnvelopeNo   string `json:"envelope_no" validate:"required"`
	RecvUsername string `json:"recv_username" validate:"required"`
	RecvUserId   string `json:"recv_user_id" validate:"required"`
	AccountNo    string `json:"account_no"`
}

type RedEnvelopeActivity struct {
	Link string `json:"link"` //活动链接
	RedEnvelopeGoodsDTO
}

func (this *RedEnvelopeActivity) CopyTo(target *RedEnvelopeActivity) {
	target.Link = this.Link
	target.EnvelopeNo = this.EnvelopeNo
	target.EnvelopeType = this.EnvelopeType
	target.Username = this.Username
	target.UserId = this.UserId
	target.Blessing = this.Blessing
	target.Amount = this.Amount
	target.AmountOne = this.AmountOne
	target.Quantity = this.Quantity
	target.RemainAmount = this.RemainAmount
	target.RemainQuantity = this.RemainQuantity
	target.ExpiredAt = this.ExpiredAt
	target.Status = this.Status
	target.OrderType = this.OrderType
	target.PayStatus = this.PayStatus
	target.CreatedAt = this.CreatedAt
	target.UpdatedAt = this.UpdatedAt
}

type RedEnvelopeGoodsDTO struct {
	EnvelopeNo     string          `json:"envelope_no"`
	EnvelopeType   int             `json:"envelope_type" validate:"required"`
	Username       string          `json:"username" validate:"required"`
	UserId         string          `json:"user_id" validate:"required"`
	Blessing       string          `json:"blessing"`
	Amount         decimal.Decimal `json:"amount" validate:"required"`
	AmountOne      decimal.Decimal `json:"amount_one"`
	Quantity       int             `json:"quantity" validate:"required"`
	RemainAmount   decimal.Decimal `json:"remain_amount"`
	RemainQuantity int             `json:"remain_quantity"`
	ExpiredAt      time.Time       `json:"expired_at"`
	Status         int             `json:"status"`
	OrderType      OrderType       `json:"order_type"`
	PayStatus      PayStatus       `json:"pay_status"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	AccountNo      string          `json:"account_no"`
}

type RedEnvelopeItemDTO struct {
	ItemNo       int64           `json:"item_no"`       // 红包订单详情编号
	EnvelopeNo   string          `json:"envelope_no"`   // 红包编号
	RecvUsername string          `json:"recv_username"` // 接收者用户名
	RecvUserId   string          `json:"recv_user_id"`  // 接收者用户id
	Amount       decimal.Decimal `json:"amount"`        // 收到金额
	Quantity     int             `json:"quantity"`      // 收到数量
	RemainAmount decimal.Decimal `json:"remain_amount"` // 剩余金额
	AccountNo    string          `json:"account_no"`    // 红包接收者账户ID
	PayStatus    int             `json:"pay_status"`    // 支付状态
	CreatedAt    time.Time       `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time       `json:"updated_at"`    // 修改时间
}
