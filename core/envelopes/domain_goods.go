package envelopes

import (
	"context"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"imooc.com/resk/infra/base"
	"imooc.com/resk/services"
	"time"
)

type goodsDomain struct {
	RedEnvelopeGoods
	item itemDomain
}

// 生成红包编号
func (domain *goodsDomain) createEnvelopeNo() {
	domain.RedEnvelopeGoods.EnvelopeNo = ksuid.New().Next().String()
}

// 创建一个红包商品对象
func (domain *goodsDomain) Create(dto *services.RedEnvelopeGoodsDTO) {
	domain.RedEnvelopeGoods.FromDTO(dto)
	domain.RemainQuantity = domain.Quantity
	domain.Username.Valid = true
	domain.Blessing.Valid = true
	domain.createEnvelopeNo()
	if domain.EnvelopeType == services.GeneralEnvelopeType {
		domain.Amount = dto.AmountOne.Mul(decimal.NewFromFloat(float64(dto.Quantity)))
	}
	if domain.EnvelopeType == services.LuckyEnvelopeType {
		domain.AmountOne = decimal.NewFromFloat(0)
	}
	domain.RemainAmount = domain.Amount
	// 过期时间
	domain.ExpiredAt = time.Now().Add(24 * time.Hour)
	domain.Status = services.OrderCreate
	domain.PayStatus = services.Paying
}

// 保存到红包商品表
func (domain *goodsDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeDao{runner: runner}
		id, err = dao.Insert(&domain.RedEnvelopeGoods)
		return err
	})
	return id, err
}

// 创建并保存红包商品
func (domain *goodsDomain) CreateAndSave(ctx context.Context, dto *services.RedEnvelopeGoodsDTO) (id int64, err error) {
	domain.Create(dto)
	return domain.Save(ctx)
}

// 查询商品信息
func (domain *goodsDomain) Get(envelopeNo string) (goods *RedEnvelopeGoods) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &RedEnvelopeDao{runner: runner}
		goods = dao.GetOne(envelopeNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return goods
}
