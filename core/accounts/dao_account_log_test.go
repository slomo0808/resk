package accounts

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/slomo0808/infra/base"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"imooc.com/resk/services"
	_ "imooc.com/resk/textx"
	"testing"
)

func TestAccountLogDao(t *testing.T) {
	Convey("通过流水编号查询", t, func() {
		err := base.Tx(func(runner *dbx.TxRunner) error {
			dao := &AccountLogDao{runner: runner}
			l := &AccountLog{
				TradeNo:         ksuid.New().Next().String(),
				LogNo:           ksuid.New().Next().String(),
				AccountNo:       ksuid.New().Next().String(),
				TargetAccountNo: ksuid.New().Next().String(),
				UserId:          ksuid.New().Next().String(),
				Username:        "账户流水测试",
				TargetUserId:    ksuid.New().Next().String(),
				TargetUsername:  "目标账户流水测试",
				Amount:          decimal.NewFromFloat(100),
				Balance:         decimal.NewFromFloat(1000),
				ChangeType:      services.AccountStoreValue,
				ChangeFlag:      services.FlagTransferIn,
				Status:          1,
				Desc:            "流水测试",
			}

			// 通过log_no查询
			Convey("通过log_no查询", func() {
				id, err := dao.Insert(l)
				So(err, ShouldBeNil)
				So(id, ShouldBeGreaterThan, 0)

				out := dao.GetOne(l.LogNo)
				So(out, ShouldNotBeNil)
				So(out.Balance.String(), ShouldEqual, l.Balance.String())
				So(out.Amount.String(), ShouldEqual, l.Amount.String())
				So(out.CreatedAt, ShouldNotBeNil)
			})

			// 通过trade_no查询
			Convey("通过trade_no查询", func() {
				id, err := dao.Insert(l)
				So(err, ShouldBeNil)
				So(id, ShouldBeGreaterThan, 0)

				out := dao.GetByTradeNo(l.TradeNo)
				So(out, ShouldNotBeNil)
				So(out.Balance.String(), ShouldEqual, l.Balance.String())
				So(out.Amount.String(), ShouldEqual, l.Amount.String())
				So(out.CreatedAt, ShouldNotBeNil)
			})

			return nil
		})
		if err != nil {
			logrus.Error(err)
		}
	})
}
