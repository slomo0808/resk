package accounts

import (
	"github.com/shopspring/decimal"
	"imooc.com/resk/services"
	"time"
)

// 账户流水表持久化对象
type AccountLog struct {
	Id              int64               `db:"id,omitempty"`
	TradeNo         string              `db:"trade_no"`
	LogNo           string              `db:"log_no,unique"`
	AccountNo       string              `db:"account_no"`
	TargetAccountNo string              `db:"target_account_no"`
	UserId          string              `db:"user_id"`
	Username        string              `db:"username"`
	TargetUserId    string              `db:"target_user_id"`
	TargetUserName  string              `db:"target_username"`
	Amount          decimal.Decimal     `db:"amount"`
	Balance         decimal.Decimal     `db:"balance"`
	ChangeType      services.ChangeType `db:"change_type"`
	ChangeFlag      services.ChangeFlag `db:"change_flag"`
	Status          int                 `db:"status"`
	Desc            string              `db:"desc"`
	CreatedAt       time.Time           `db:"created_at,omitempty"`
}
