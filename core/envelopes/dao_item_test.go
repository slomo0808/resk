package envelopes

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"imooc.com/resk/infra/base"
	"imooc.com/resk/services"
	_ "imooc.com/resk/textx"
	"testing"
)

func TestRedEnvelopeItemDao(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		Convey("红包item数据库操作", t, func() {
			data := &RedEnvelopeItem{
				ItemNo:       1234567890,
				EnvelopeNo:   ksuid.New().Next().String(),
				RecvUsername: "测试用item Username",
				RecvUserId:   "测试用item UserId",
				Amount:       decimal.NewFromFloat(10.2),
				Quantity:     1,
				RemainAmount: decimal.NewFromFloat(100),
				AccountNo:    "测试用item AccountNo",
				PayStatus:    int(services.Payed),
			}
			// 插入操作
			id, err := dao.Insert(data)
			So(id, ShouldBeGreaterThan, 0)
			So(err, ShouldBeNil)

			// 查询操作
			out := dao.GetOne(data.ItemNo)
			So(out, ShouldNotBeNil)
			So(out.EnvelopeNo, ShouldEqual, data.EnvelopeNo)
			So(out.RecvUserId, ShouldEqual, data.RecvUserId)
			So(out.RecvUsername, ShouldEqual, data.RecvUsername)
			So(out.Amount.String(), ShouldEqual, data.Amount.String())
			So(out.RemainAmount.String(), ShouldEqual, data.RemainAmount.String())
			So(out.AccountNo, ShouldEqual, data.AccountNo)
			So(out.PayStatus, ShouldEqual, data.PayStatus)
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}
