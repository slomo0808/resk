package accounts

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/slomo0808/infra/base"
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
	if err := base.ValidateStruct(dto); err != nil {
		return nil, err
	}
	// 验证账户是否已经存在
	acc := domain.GetAccountByUserIdAndType(dto.UserId, services.EnvelopeAccountType)
	if acc != nil {
		return acc, errors.New(fmt.Sprintf("用户的该类型账户已经存在，username=%s[%s],账户类型:%d",
			acc.Username, acc.UserId, acc.AccountType))
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
	if err := base.ValidateStruct(dto); err != nil {
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
