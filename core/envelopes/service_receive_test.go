package envelopes

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"imooc.com/resk/services"
	"strconv"
	"testing"
)

func TestGoodsDomain_Receive(t *testing.T) {
	// 1.准备几个h红包资金账户，用户收发红包
	as := services.GetAccountService()
	accounts := make([]*services.AccountDTO, 0)
	size := 10
	Convey("收红包测试用例", t, func() {
		for i := 0; i < size; i++ {
			account := &services.AccountCreatedDTO{
				UserId:      ksuid.New().Next().String(),
				Username:    "测试用户" + strconv.Itoa(i),
				AccountName: "测试账户" + strconv.Itoa(i),
				AccountType: int(services.EnvelopeAccountType),
				Amount:      "10000",
			}
			// 账户创建
			acDto, err := as.CreateAccount(account)
			So(err, ShouldBeNil)
			So(acDto, ShouldNotBeNil)
			accounts = append(accounts, acDto)
		}
		acDTO := accounts[0]
		res := services.GetRedEnvelopeService()
		// 2.使用其中一个发红包
		// 发送普通红包
		goods := &services.RedEnvelopeSendingDTO{
			EnvelopeType: services.GeneralEnvelopeType,
			Username:     acDTO.Username,
			UserId:       acDTO.UserId,
			Amount:       decimal.NewFromFloat(10),
			Quantity:     size,
			Blessing:     "收红包测试",
		}
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

		// 3.使用发送红包数量的人收红包
		remainAmount := at.Amount
		Convey("收普通红包", func() {
			for _, account := range accounts {
				rcv := &services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   at.EnvelopeNo,
					RecvUsername: account.Username,
					RecvUserId:   account.UserId,
					AccountNo:    account.AccountNo,
				}
				item, err := res.Receive(rcv)
				So(err, ShouldBeNil)
				So(item, ShouldNotBeNil)
				So(item.Amount.String(), ShouldEqual, at.AmountOne.String())
				remainAmount = remainAmount.Sub(item.Amount)
				So(item.RemainAmount.String(), ShouldEqual, remainAmount.String())
			}
		})
		Convey("收运气红包", func() {
			goods.EnvelopeType = services.LuckyEnvelopeType
			at, err := res.SendOut(*goods)
			So(at, ShouldNotBeNil)
			So(err, ShouldBeNil)
			remainAmount = at.RemainAmount
			accounts = accounts[:10]
			for _, account := range accounts {
				rcv := &services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   at.EnvelopeNo,
					RecvUsername: account.Username,
					RecvUserId:   account.UserId,
					AccountNo:    account.AccountNo,
				}
				item, err := res.Receive(rcv)

				So(err, ShouldBeNil)
				So(item, ShouldNotBeNil)
				remainAmount = remainAmount.Sub(item.Amount)
				So(item.RemainAmount.String(), ShouldEqual, remainAmount.String())
			}
		})
	})
}
