package envelopes

import (
	"context"
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/tietang/dbx"
	"imooc.com/resk/core/accounts"
	"imooc.com/resk/infra/algo"
	"imooc.com/resk/infra/base"
	"imooc.com/resk/services"
)

var multiple = decimal.NewFromFloat(100.0)

// 收红包业务逻辑代码
func (domain *goodsDomain) Receive(
	ctx context.Context,
	dto *services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	// 1.创建收红包的订单明细
	domain.preCreateItem(dto)
	// 2.查询出当前红包的剩余数量和剩余金额信息
	goods := domain.Get(dto.EnvelopeNo)
	// 3.校验剩余红包数量和剩余金额
	// 如果没有剩余，直接返回无可用红包金额
	if goods.RemainQuantity <= 0 || goods.RemainAmount.Cmp(decimal.NewFromFloat(0)) <= 0 {
		return nil, errors.New("没有足够的金额")
	}
	// 4.使用红包算法计算红包金额
	amount := domain.nextAmount(goods)

	err = base.Tx(func(runner *dbx.TxRunner) error {
		// 5.使用乐观锁更新语句，尝试更新剩余数量和剩余金额
		dao := RedEnvelopeDao{runner: runner}
		rows, err := dao.UpdateBalance(goods.EnvelopeNo, amount)
		// 如果更新失败，row effect返回0，表示无可用红包数量与金额
		if rows <= 0 || err != nil {
			return errors.New("没有足够的金额")
		}
		// 如果更新成功，row effect返回1，表示抢到红包
		// 6.保存订单明细数据
		domain.item.Quantity = 1
		domain.item.PayStatus = int(services.Paying)
		domain.item.AccountNo = dto.AccountNo
		domain.item.RemainAmount = goods.RemainAmount.Sub(amount)
		domain.item.Amount = amount
		txCtx := base.WithValueContext(ctx, runner)
		_, err = domain.item.Save(txCtx)
		if err != nil {
			return err
		}
		// 7.将抢到的红包金额从系统红包中间账户转入当前用户的资金账户
		status, err := domain.transfer(txCtx, dto)
		if status == services.TransferredStatusSuccess {
			return nil
		}
		return err
	})
	return domain.item.ToDTO(), err
}

func (domain *goodsDomain) transfer(ctx context.Context,
	dto *services.RedEnvelopeReceiveDTO) (status services.TransferredStatus, err error) {
	systemAccount := base.GetSystemAccount()
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	target := services.TradeParticipator{
		AccountNo: dto.AccountNo,
		UserId:    dto.RecvUserId,
		Username:  dto.RecvUsername,
	}
	if target.AccountNo == "" {
		ac := accounts.NewAccountDomain().GetAccountByUserIdAndType(target.UserId, services.EnvelopeAccountType)
		target.AccountNo = ac.AccountNo
	}
	transferDTO := &services.AccountTransferDTO{
		TradeBody:   body,
		TradeTarget: target,
		Amount:      domain.item.Amount,
		ChangeType:  services.EnvelopeOutgoing,
		ChangeFlag:  services.FlagTransferOut,
		Desc:        "红包金额扣减",
	}
	accountDomain := accounts.NewAccountDomain()
	// 系统账户扣减资金
	status, err = accountDomain.TransferWithContextTx(ctx, transferDTO)
	if status != services.TransferredStatusSuccess {
		return
	}
	// 用户账户增加资金
	transferDTO = &services.AccountTransferDTO{
		TradeBody:   target,
		TradeTarget: body,
		Amount:      domain.item.Amount,
		ChangeType:  services.EnvelopeIncoming,
		ChangeFlag:  services.FlagTransferIn,
		Desc:        "抢红包成功收入",
	}
	return accountDomain.TransferWithContextTx(ctx, transferDTO)
}

// 创建收红包的订单明细
func (domain *goodsDomain) preCreateItem(dto *services.RedEnvelopeReceiveDTO) {
	domain.item.AccountNo = dto.AccountNo
	domain.item.EnvelopeNo = dto.EnvelopeNo
	domain.item.RecvUsername = sql.NullString{
		String: dto.RecvUsername,
		Valid:  true,
	}
	domain.item.RecvUserId = dto.RecvUserId
	domain.item.createItemNo()
}

// 计算红包金额
func (domain *goodsDomain) nextAmount(goods *RedEnvelopeGoods) (amount decimal.Decimal) {
	if goods.RemainQuantity == 1 {
		return goods.RemainAmount
	}
	if goods.EnvelopeType == services.GeneralEnvelopeType {
		return goods.AmountOne
	}
	cent := goods.RemainAmount.Mul(multiple).IntPart()
	next := algo.DoubleAverage(int64(goods.RemainQuantity), cent)
	amount = decimal.NewFromFloat(float64(next)).Div(multiple)
	return amount
}
