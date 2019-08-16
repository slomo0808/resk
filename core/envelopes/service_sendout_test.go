package envelopes

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"imooc.com/resk/infra/base"
	"imooc.com/resk/services"
	"testing"
)

func TestRedEnvelopeService_SendOut(t *testing.T) {
	// 发红包人的资金账户
	as := services.GetAccountService()
	acDTO := &services.AccountCreatedDTO{
		UserId:      ksuid.New().Next().String(),
		Username:    "测试用户",
		AccountName: "测试账户",
		AccountType: int(services.EnvelopeAccountType),
		Amount:      "1000",
	}
	res := services.GetRedEnvelopeService()

	err := base.Tx(func(runner *dbx.TxRunner) error {
		Convey("准备账户", t, func() {
			aDTO, err := as.CreateAccount(acDTO)
			So(aDTO, ShouldNotBeNil)
			So(err, ShouldBeNil)
			So(aDTO.Balance.String(), ShouldEqual, acDTO.Amount)
		})

		Convey("发红包测试代码", t, func() {
			goods := &services.RedEnvelopeSendingDTO{
				EnvelopeType: services.GeneralEnvelopeType,
				Username:     acDTO.Username,
				UserId:       acDTO.UserId,
				Amount:       decimal.NewFromFloat(10),
				Quantity:     3,
			}

			Convey("普通红包", func() {
				at, err := res.SendOut(*goods)
				So(at, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(at.Link, ShouldNotBeEmpty)
				So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)
				// 验证每一个属性
				dto := at.RedEnvelopeGoodsDTO
				So(dto.Username, ShouldEqual, goods.Username)
				So(dto.UserId, ShouldEqual, goods.UserId)
				So(dto.Amount.String(), ShouldEqual,
					goods.Amount.Mul(decimal.NewFromFloat(float64(goods.Quantity))).String())
				So(dto.AmountOne.String(), ShouldEqual, goods.Amount.String())
			})

			Convey("碰运气红包", func() {
				goods.EnvelopeType = services.LuckyEnvelopeType
				at, err := res.SendOut(*goods)
				So(at, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(at.Link, ShouldNotBeEmpty)
				So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)
				// 验证每一个属性
				dto := at.RedEnvelopeGoodsDTO
				So(dto.Username, ShouldEqual, goods.Username)
				So(dto.UserId, ShouldEqual, goods.UserId)
				So(dto.Amount.String(), ShouldEqual, goods.Amount.String())
				So(dto.AmountOne.String(), ShouldEqual, decimal.NewFromFloat(0).String())
			})
		})

		return nil
	})
	if err != nil {
		logrus.Panic(err)
	}
}
