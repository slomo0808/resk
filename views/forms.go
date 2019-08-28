package views

import "github.com/shopspring/decimal"

type UserForm struct {
	UserId   string `form:"user_id"`
	Mobile   string `form:"mobile"`
	Username string `form:"username"`
}

type RedEnvelopeSendingFrom struct {
	EnvelopeType int             `form:"envelopeType"`
	Blessing     string          `form:"blessing"`
	Amount       decimal.Decimal `form:"amount"`
	Quantity     int             `form:"quantity"`
}
