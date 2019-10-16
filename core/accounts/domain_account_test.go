package accounts

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/slomo0808/infra/base"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"imooc.com/resk/services"
	"testing"
)

func TestAccountDomain_Create(t *testing.T) {
	dto := &services.AccountDTO{
		AccountName: "账户创建测试",
		UserId:      ksuid.New().Next().String(),
		Username:    "账户创建测试",
		Balance:     decimal.NewFromFloat(0),
		Status:      1,
	}
	domain := accountDomain{}
	Convey("账户创建测试", t, func() {
		adto, err := domain.Create(dto)
		So(err, ShouldBeNil)
		So(adto, ShouldNotBeNil)
		So(adto.Balance.String(), ShouldEqual, dto.Balance.String())
		So(adto.Username, ShouldEqual, dto.Username)
		So(adto.UserId, ShouldEqual, dto.UserId)
		So(adto.Status, ShouldEqual, dto.Status)
		So(adto.AccountName, ShouldEqual, dto.AccountName)

	})
}

func TestAccountDomain_Transfer(t *testing.T) {
	// 两个账户，一个交易主体，一个交易对象
	body := &services.AccountDTO{
		UserId:      ksuid.New().Next().String(),
		Username:    "交易主体",
		Balance:     decimal.NewFromFloat(1000),
		Status:      1,
		AccountType: int(services.EnvelopeAccountType),
	}
	target := &services.AccountDTO{
		UserId:      ksuid.New().Next().String(),
		Username:    "交易对象",
		Balance:     decimal.NewFromFloat(0),
		Status:      1,
		AccountType: int(services.EnvelopeAccountType),
	}
	domain := accountDomain{}
	Convey("转账测试", t, func() {
		// 交易主体创建
		aBody, err := domain.Create(body)
		So(err, ShouldBeNil)
		So(aBody, ShouldNotBeNil)
		So(aBody.Balance.String(), ShouldEqual, body.Balance.String())
		So(aBody.Username, ShouldEqual, body.Username)
		So(aBody.UserId, ShouldEqual, body.UserId)
		So(aBody.Status, ShouldEqual, body.Status)
		So(aBody.AccountName, ShouldEqual, body.AccountName)

		// 交易目标创建
		aTarget, err := domain.Create(target)
		So(err, ShouldBeNil)
		So(aTarget, ShouldNotBeNil)
		So(aTarget.Balance.String(), ShouldEqual, target.Balance.String())
		So(aTarget.Username, ShouldEqual, target.Username)
		So(aTarget.UserId, ShouldEqual, target.UserId)
		So(aTarget.Status, ShouldEqual, target.Status)
		So(aTarget.AccountName, ShouldEqual, target.AccountName)

		// 转账操作验证
		Convey("余额充足，转入其他账户", func() {
			dto := &services.AccountTransferDTO{
				TradeNo: ksuid.New().Next().String(),
				TradeBody: services.TradeParticipator{
					AccountNo: aBody.AccountNo,
					UserId:    aBody.UserId,
					Username:  aBody.Username,
				},
				TradeTarget: services.TradeParticipator{
					AccountNo: aTarget.AccountNo,
					UserId:    aTarget.UserId,
					Username:  aTarget.Username,
				},
				AmountStr:  "100",
				Amount:     decimal.NewFromFloat(100),
				ChangeType: services.EnvelopeOutgoing,
				ChangeFlag: services.FlagTransferOut,
				Desc:       "余额充足，转入其他账户",
			}
			scode, err := domain.Transfer(dto)
			So(err, ShouldBeNil)
			So(scode, ShouldEqual, services.TransferredStatusSuccess)

			// 实际余额更新后的值
			err = base.Tx(func(runner *dbx.TxRunner) error {
				accountDao := AccountDao{runner: runner}
				accountLogDao := AccountLogDao{runner: runner}
				// 得到交易主体新数据
				bodyOut := accountDao.GetOne(aBody.AccountNo)
				bodyOutByUserId := domain.GetEnvelopeAccountByUserId(aBody.UserId)
				So(bodyOut.UserId, ShouldEqual, bodyOutByUserId.UserId)
				So(bodyOut.AccountNo, ShouldEqual, bodyOutByUserId.AccountNo)
				So(bodyOut.Status, ShouldEqual, bodyOutByUserId.Status)
				So(bodyOut.AccountType, ShouldEqual, bodyOutByUserId.AccountType)
				So(bodyOut.Balance.String(), ShouldEqual, bodyOutByUserId.Balance.String())

				So(bodyOut, ShouldNotBeNil)
				// 得到交易对象最新数据
				targetOut := accountDao.GetOne(aTarget.AccountNo)
				targetOutByUserId := domain.GetEnvelopeAccountByUserId(aTarget.UserId)
				So(targetOut, ShouldNotBeNil)
				So(targetOut.UserId, ShouldEqual, targetOutByUserId.UserId)
				So(targetOut.AccountNo, ShouldEqual, targetOutByUserId.AccountNo)
				So(targetOut.Status, ShouldEqual, targetOutByUserId.Status)
				So(targetOut.AccountType, ShouldEqual, targetOutByUserId.AccountType)
				So(targetOut.Balance.String(), ShouldEqual, targetOutByUserId.Balance.String())

				// 验证交易主体，余额扣减正确
				So(bodyOut.Balance.String(),
					ShouldEqual,
					body.Balance.Add(dto.Amount.Mul(decimal.NewFromFloat(-1))).String())
				// 验证交易对象，余额增加正确
				// So(targetOut.Balance.String(), ShouldEqual, target.Balance.Add(dto.Amount).String())
				// 验证交易流水数据正确
				outLog := accountLogDao.GetByTradeNo(dto.TradeNo)

				So(outLog, ShouldNotBeNil)
				// 验证流水中交易主体和交易目标正确
				So(outLog.AccountNo, ShouldEqual, bodyOut.AccountNo)
				So(outLog.UserId, ShouldEqual, bodyOut.UserId)
				So(outLog.Username, ShouldEqual, bodyOut.Username.String)
				So(outLog.TargetAccountNo, ShouldEqual, targetOut.AccountNo)
				So(outLog.TargetUserId, ShouldEqual, targetOut.UserId)
				So(outLog.TargetUsername, ShouldEqual, targetOut.Username.String)
				// 验证交易流水中交易主体的余额正确
				So(outLog.Balance.String(), ShouldEqual, bodyOut.Balance.String())
				// 验证交易流水中交易的金额正确
				So(outLog.Amount.String(), ShouldEqual, dto.Amount.String())
				// 验证交易流水changeFlag和changeType正确
				So(outLog.ChangeFlag, ShouldEqual, dto.ChangeFlag)
				So(outLog.ChangeType, ShouldEqual, dto.ChangeType)

				return nil
			})
			So(err, ShouldBeNil)
		})
		Convey("余额不足，金额转出", func() {
			dto := &services.AccountTransferDTO{
				TradeNo: ksuid.New().Next().String(),
				TradeBody: services.TradeParticipator{
					AccountNo: aBody.AccountNo,
					UserId:    aBody.UserId,
					Username:  aBody.Username,
				},
				TradeTarget: services.TradeParticipator{
					AccountNo: aTarget.AccountNo,
					UserId:    aTarget.UserId,
					Username:  aTarget.Username,
				},
				AmountStr:  "100000",
				Amount:     decimal.NewFromFloat(100000),
				ChangeType: services.EnvelopeOutgoing,
				ChangeFlag: services.FlagTransferOut,
				Desc:       "余额不足，转入其他账户",
			}
			scode, err := domain.Transfer(dto)
			So(err, ShouldNotBeNil)
			So(scode, ShouldEqual, services.TransferredStatusSufficientFunds)
			// 实际余额更新后的值
			err = base.Tx(func(runner *dbx.TxRunner) error {
				accountDao := AccountDao{runner: runner}
				accountLogDao := AccountLogDao{runner: runner}
				// 得到交易主体新数据
				bodyOut := accountDao.GetOne(aBody.AccountNo)
				So(bodyOut, ShouldNotBeNil)
				// 得到交易对象最新数据
				targetOut := accountDao.GetOne(aTarget.AccountNo)
				So(targetOut, ShouldNotBeNil)

				// 验证交易主体，余额扣减正确,因为扣减失败所以没有变化
				So(bodyOut.Balance.String(), ShouldEqual, body.Balance.String())
				// 验证交易对象，余额增加正确, 交易失败，没有增加
				// So(targetOut.Balance.String(), ShouldEqual, target.Balance.String())
				// 验证交易流水数据正确, 因为没有形成交易，所以无流水
				outLog := accountLogDao.GetByTradeNo(dto.TradeNo)
				So(outLog, ShouldBeNil)
				return nil
			})
			So(err, ShouldBeNil)
		})
		Convey("充值", func() {
			dto := &services.AccountTransferDTO{
				TradeNo: ksuid.New().Next().String(),
				TradeBody: services.TradeParticipator{
					AccountNo: aBody.AccountNo,
					UserId:    aBody.UserId,
					Username:  aBody.Username,
				},
				TradeTarget: services.TradeParticipator{
					AccountNo: aBody.AccountNo,
					UserId:    aBody.UserId,
					Username:  aBody.Username,
				},
				AmountStr:  "100",
				Amount:     decimal.NewFromFloat(100),
				ChangeType: services.AccountStoreValue,
				ChangeFlag: services.FlagTransferIn,
				Desc:       "充值",
			}
			scode, err := domain.Transfer(dto)
			So(err, ShouldBeNil)
			So(scode, ShouldEqual, services.TransferredStatusSuccess)
			// 实际余额更新后的值
			err = base.Tx(func(runner *dbx.TxRunner) error {
				accountDao := AccountDao{runner: runner}
				accountLogDao := AccountLogDao{runner: runner}
				// 得到交易主体新数据
				bodyOut := accountDao.GetOne(aBody.AccountNo)
				So(bodyOut, ShouldNotBeNil)

				// 充值，验证余额增加正确
				So(bodyOut.Balance.String(),
					ShouldEqual,
					body.Balance.Add(dto.Amount).String())

				// 验证交易流水数据正确
				outLog := accountLogDao.GetByTradeNo(dto.TradeNo)
				So(outLog, ShouldNotBeNil)
				// 验证流水中交易主体和交易目标正确
				So(outLog.AccountNo, ShouldEqual, bodyOut.AccountNo)
				So(outLog.UserId, ShouldEqual, bodyOut.UserId)
				So(outLog.Username, ShouldEqual, bodyOut.Username.String)
				So(outLog.TargetAccountNo, ShouldEqual, bodyOut.AccountNo)
				So(outLog.TargetUserId, ShouldEqual, bodyOut.UserId)
				So(outLog.TargetUsername, ShouldEqual, bodyOut.Username.String)
				// 验证交易流水中交易主体的余额正确
				So(outLog.Balance.String(), ShouldEqual, bodyOut.Balance.String())
				// 验证交易流水中交易的金额正确
				So(outLog.Amount.String(), ShouldEqual, dto.Amount.String())
				// 验证交易流水changeFlag和changeType正确
				So(outLog.ChangeFlag, ShouldEqual, dto.ChangeFlag)
				So(outLog.ChangeType, ShouldEqual, dto.ChangeType)
				return nil
			})
			So(err, ShouldBeNil)
		})
	})
}
