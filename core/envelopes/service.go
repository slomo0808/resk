package envelopes

import (
	"errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"imooc.com/resk/infra/base"
	"imooc.com/resk/services"
	"sync"
)

var once sync.Once

func init() {
	once.Do(func() {
		services.IRedEnvelopeService = new(redEnvelopeService)
	})
}

type redEnvelopeService struct{}

func (s *redEnvelopeService) SendOut(dto services.RedEnvelopeSendingDTO) (*services.RedEnvelopeActivity, error) {
	// 验证
	if err := base.ValidateStruct(&dto); err != nil {
		return nil, err
	}
	// 获取红包发送人的资金信息
	account := services.GetAccountService().GetEnvelopeAccountByUserId(dto.UserId)
	if account == nil {
		return nil, errors.New("用户的账户不存在:" + dto.UserId)
	}
	goods := (&dto).ToGoods()
	goods.AccountNo = account.AccountNo

	if goods.Blessing == "" {
		goods.Blessing = services.DefaultBlessing
	}
	if goods.EnvelopeType == services.GeneralEnvelopeType {
		goods.AmountOne = goods.Amount
		goods.Amount = decimal.Decimal{}
	}
	// 执行发红包的逻辑
	domain := new(goodsDomain)
	activity, err := domain.SendOut(goods)
	if err != nil {
		logrus.Error(err)
	}
	return activity, err
}

func (s *redEnvelopeService) Receive(*services.RedEnvelopeReceiveDTO) (*services.RedEnvelopeItemDTO, error) {
	panic("implement me")
}

func (s *redEnvelopeService) Refund(string) *services.RedEnvelopeGoodsDTO {
	panic("implement me")
}

func (s *redEnvelopeService) Get() services.RedEnvelopeGoodsDTO {
	panic("implement me")
}
