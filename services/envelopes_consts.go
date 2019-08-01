package services

const (
	DefaultBlessing = "恭喜发财"
)

// 订单类型 ： 发布单，退款单
type OrderType int

const (
	OrderTypeSending OrderType = 1
	OrderTypeRefund  OrderType = 2
)

// 支付状态：未支付，支付中，已支付，支付失败
type PayStatus int

const (
	PayNothing PayStatus = 1
	Paying     PayStatus = 2
	Payed      PayStatus = 3
	PayFailed  PayStatus = 4
)

// 红包订单状态：创建， 发布，过期，失效
type OrderStatus int

const (
	OrderCreate   OrderStatus = 1
	OrderSending  OrderStatus = 2
	OrderExpired  OrderStatus = 3
	OrderDisabled OrderStatus = 4
)

// 活动状态，创建，激活，过期，失效
//type ActivityStatus int
//
//const (
//	ActivityCreate   ActivityStatus = 1
//	ActivitySending  ActivityStatus = 2
//	ActivityExpired  ActivityStatus = 3
//	ActivityDisabled ActivityStatus = 4
//)

// 红包类型
type EnvelopeType int

const (
	GeneralEnvelopeType = 1
	LuckyEnvelopeType   = 2
)