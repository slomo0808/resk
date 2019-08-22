package envelopes

import (
	"context"
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"imooc.com/resk/infra/base"
	"imooc.com/resk/services"
)

const (
	pageSize = 100
)

type ExpiredEnvelopeDomain struct {
	expiredGoods []RedEnvelopeGoods
	offset       int
}

// 查询出过期红包
func (d *ExpiredEnvelopeDomain) Next() (ok bool) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeDao{runner: runner}
		d.expiredGoods = dao.FindExpired(d.offset, pageSize)
		logrus.Infof("查询到%d个可退款红包", len(d.expiredGoods))
		if len(d.expiredGoods) > 0 {
			d.offset += len(d.expiredGoods)

			ok = true
		}
		return nil
	})
	if err != nil {
		return false
	}
	return ok
}

func (d *ExpiredEnvelopeDomain) Expired() (err error) {
	for d.Next() {
		for _, g := range d.expiredGoods {
			logrus.Debugf("过期红包退款开始: %+v", g)
			err = d.ExpiredOne(g)
			if err != nil {
				logrus.Error(err)
			}
			logrus.Debugf("过期红包退款结束: %+v", g)
		}
	}
	return err
}

// 发起一个退款流程
func (d *ExpiredEnvelopeDomain) ExpiredOne(goods RedEnvelopeGoods) (err error) {
	// 创建一个退款订单
	refund := goods
	refund.OrderType = services.OrderTypeRefund
	refund.RemainAmount = decimal.NewFromFloat(0)
	refund.RemainQuantity = 0
	refund.Status = services.OrderExpired
	refund.PayStatus = services.Refunding
	refund.OriginEnvelopeNo = sql.NullString{
		String: goods.EnvelopeNo,
		Valid:  true,
	}
	refund.EnvelopeNo = ""
	domain := goodsDomain{RedEnvelopeGoods: refund}
	domain.createEnvelopeNo()
	err = base.Tx(func(runner *dbx.TxRunner) error {
		txCtx := base.WithValueContext(context.Background(), runner)
		id, err := domain.Save(txCtx)
		if err != nil || id <= 0 {
			return errors.New("创建退款订单失败")
		}

		// 修改原订单状态
		dao := RedEnvelopeDao{runner: runner}
		rows, err := dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpired)
		if err != nil {
			return errors.New("更新原订单为过期状态-失败" + err.Error())
		}
		if rows == 0 {
			return errors.New("更新原订单为过期状态-失败, 影响行数为0")
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 调用资金账户接口转账
	systemAccount := base.GetSystemAccount()
	account := services.GetAccountService().GetEnvelopeAccountByUserId(goods.UserId)
	if account == nil {
		return errors.New("没有找到该用户的红包资金账户")
	}
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	target := services.TradeParticipator{
		AccountNo: account.AccountNo,
		UserId:    account.UserId,
		Username:  account.Username,
	}
	// 系统账户扣减资金
	transfer := &services.AccountTransferDTO{
		TradeNo:     domain.RedEnvelopeGoods.EnvelopeNo,
		TradeBody:   body,
		TradeTarget: target,
		Amount:      goods.RemainAmount,
		AmountStr:   goods.RemainAmount.String(),
		ChangeType:  services.SysEnvelopeExpiredRefund,
		ChangeFlag:  services.FlagTransferOut,
		Desc:        "过期退款，系统账户扣减资金：" + goods.EnvelopeNo,
	}
	status, err := services.GetAccountService().Transfer(transfer)
	if status != services.TransferredStatusSuccess {
		return err
	}
	// 用户账户增加资金
	transfer = &services.AccountTransferDTO{
		TradeNo:     domain.RedEnvelopeGoods.EnvelopeNo,
		TradeBody:   target,
		TradeTarget: body,
		Amount:      goods.RemainAmount,
		AmountStr:   goods.RemainAmount.String(),
		ChangeType:  services.EnvelopeExpiredRefund,
		ChangeFlag:  services.FlagTransferIn,
		Desc:        "过期退款，返还用户资金：" + goods.EnvelopeNo,
	}
	status, err = services.GetAccountService().Transfer(transfer)
	if status != services.TransferredStatusSuccess {
		return err
	}

	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeDao{runner: runner}
		// 修改原订单状态
		rows, err := dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpiredRefundSucceed)
		if err != nil {
			return errors.New("更新原订单状态为退款成功状态-失败" + err.Error())
		}
		if rows == 0 {
			return errors.New("更新原订单状态为退款成功状态, 影响行数为0")
		}
		// 修改退款订单状态
		rows, err = dao.UpdateOrderStatus(refund.EnvelopeNo, services.OrderExpiredRefundSucceed)
		if err != nil {
			return errors.New("更新原订单为退款成功状态-失败" + err.Error())
		}
		if rows == 0 {
			return errors.New("更新原订单为退款成功状态-失败, 影响行数为0")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
