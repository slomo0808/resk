package envelopes

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	accountServices "github.com/slomo0808/account/services"
	"github.com/slomo0808/infra/base"
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
	account := accountServices.GetAccountService().GetEnvelopeAccountByUserId(dto.UserId)
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

func (s *redEnvelopeService) Receive(dto *services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	// 参数校验
	if err = base.ValidateStruct(dto); err != nil {
		return nil, err
	}
	// 获取当前收红包用户的账户信息
	account := accountServices.GetAccountService().GetEnvelopeAccountByUserId(dto.RecvUserId)
	if account == nil {
		return nil, errors.New("红包资金账户不存在:user_id = " + dto.RecvUserId)
	}
	// 进行尝试收红包
	domain := new(goodsDomain)
	item, err = domain.Receive(context.Background(), dto)
	return item, err
}

func (s *redEnvelopeService) Refund(string) *services.RedEnvelopeGoodsDTO {
	panic("implement me")
}

func (s *redEnvelopeService) Get(envelopeNo string) *services.RedEnvelopeGoodsDTO {
	domain := new(goodsDomain)
	goods := domain.Get(envelopeNo)
	return goods.ToDTO()
}

func (s *redEnvelopeService) ListSent(userId string, page, size int) []*services.RedEnvelopeGoodsDTO {
	domain := new(goodsDomain)
	pos := domain.FindByUser(userId, page, size)

	orders := make([]*services.RedEnvelopeGoodsDTO, 0, len(pos))
	for _, po := range pos {
		orders = append(orders, po.ToDTO())
	}
	return orders
}

func (s *redEnvelopeService) ListReceived(userId string, page, size int) (items []*services.RedEnvelopeItemDTO) {
	domain := new(goodsDomain)
	pos := domain.ListReceived(userId, page, size)
	items = make([]*services.RedEnvelopeItemDTO, 0, len(pos))
	if len(pos) == 0 {
		return items
	}
	for _, p := range pos {
		items = append(items, p.ToDTO())
	}
	return
}

func (s *redEnvelopeService) ListItems(envelopeNo string) (items []*services.RedEnvelopeItemDTO) {
	domain := itemDomain{}
	return domain.FindItems(envelopeNo)
}

func (s *redEnvelopeService) ListReceivable(offset int, size int) []*services.RedEnvelopeGoodsDTO {
	domain := new(goodsDomain)
	pos := domain.ListReceivable(offset, size)
	orders := make([]*services.RedEnvelopeGoodsDTO, 0, len(pos))
	for _, po := range pos {
		orders = append(orders, po.ToDTO())
	}
	return orders
}
