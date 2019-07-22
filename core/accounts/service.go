package accounts

import (
	"errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	"imooc.com/resk/infra/base"
	"imooc.com/resk/services"
	"sync"
)

var _ services.AccountService = new(accountService)
var one sync.Once

func init() {
	one.Do(func() {
		services.IAccountService = new(accountService)
	})
}

type accountService struct {
}

func (s *accountService) CreateAccount(dto *services.AccountCreatedDTO) (*services.AccountDTO, error) {
	domain := NewAccountDomain()
	// 验证输入参数
	err := base.Validate().Struct(dto)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("账户创建验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				logrus.Error(e.Translate(base.Translate()))
			}
		}
		return nil, err
	}
	// 执行账户创建的业务逻辑代码
	amount, err := decimal.NewFromString(dto.Amount)
	if err != nil {
		return nil, err
	}
	account := &services.AccountDTO{
		UserId:       dto.UserId,
		Username:     dto.Username,
		AccountName:  dto.AccountName,
		AccountType:  dto.AccountType,
		CurrencyCode: dto.CurrencyCode,
		Status:       1,
		Balance:      amount,
	}
	return domain.Create(account)
}

func (s *accountService) Transfer(dto *services.AccountTransferDTO) (services.TransferredStatus, error) {
	domain := NewAccountDomain()
	// 验证dto
	err := base.Validate().Struct(dto)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logrus.Error("转账验证错误", err)
		}
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				logrus.Error(e.Translate(base.Translate()))
			}
		}
		return services.TransferredStatusFailure, err
	}
	// 执行转账逻辑
	amount, err := decimal.NewFromString(dto.AmountStr)
	if err != nil {
		return services.TransferredStatusFailure, err
	}
	dto.Amount = amount
	if dto.ChangeFlag == services.FlagTransferOut {
		if dto.ChangeType > 0 {
			return services.TransferredStatusFailure,
				errors.New("如果changeFlag为支出，那么changeType必须小于0")
		}
	} else {
		if dto.ChangeType < 0 {
			return services.TransferredStatusFailure,
				errors.New("如果changeFlag为收入，那么changeType必须大于0")
		}
	}
	return domain.Transfer(dto)
}

func (s *accountService) StoreValue(dto *services.AccountTransferDTO) (services.TransferredStatus, error) {
	dto.TradeTarget = dto.TradeBody
	dto.ChangeFlag = services.FlagTransferIn
	dto.ChangeType = services.AccountStoreValue
	return s.Transfer(dto)
}

func (s *accountService) GetAccount(accountNo string) *services.AccountDTO {
	domain := NewAccountDomain()
	return domain.GetAccount(accountNo)
}

func (s *accountService) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
	domain := NewAccountDomain()
	return domain.GetEnvelopeAccountByUserId(userId)
}
