package services

import "time"

type AccountService interface {
	// 创建账户
	CreateAccount(dto AccountCreatedDTO) (*AccountDTO, error)
	// 转账
	Transfer(dto AccountTransferDTO) (TransferredStatus, error)
	// 储值
	StoreValue(dto AccountTransferDTO) (TransferredStatus, error)
	// 查询
	GetEnvelopeAccountByUserId(userId string) *AccountDTO
}

// DTO : Data Transfer Object
// 数据传输对象

// 账户创建
type AccountCreatedDTO struct {
	UserId       string
	Username     string
	AccountName  string
	AccountType  int
	CurrencyCode string
	Amount       string
}

// 账户信息
type AccountDTO struct {
	AccountCreatedDTO
	AccountNo string
	CreatedAt time.Time
}

// 账户交易参与者
type TradeParticipator struct {
	AccountNo string
	UserId    string
	Username  string
}

// 账户转账
type AccountTransferDTO struct {
	TradeNo     string
	TradeBody   TradeParticipator
	TradeTarget TradeParticipator
	AmountStr   string
	ChangeType  ChangeType
	ChangeFlag  ChangeFlag
	Desc        string
}
