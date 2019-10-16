package envelopes

import (
	"context"
	"errors"
	"github.com/slomo0808/account/core/accounts"
	accountServices "github.com/slomo0808/account/services"
	"github.com/slomo0808/infra/base"
	"github.com/tietang/dbx"
	"imooc.com/resk/services"
	"path"
)

// 发红包业务领域代码
func (domain *goodsDomain) SendOut(
	dto *services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 创建红包商品
	domain.Create(dto)
	// 创建活动
	activity = &services.RedEnvelopeActivity{}
	// 红包链接 格式 http://域名/v1/envelope/{id}/link/
	link := base.GetEnvelopeActivityLink()
	d := base.GetEnvelopeDomain()
	activity.Link = path.Join(d, link, domain.EnvelopeNo)

	err = base.Tx(func(runner *dbx.TxRunner) error {
		//
		ctx := base.WithValueContext(context.Background(), runner)
		// 保存红包商品
		id, err := domain.Save(ctx)
		if id < 1 {
			return errors.New("id less then 1, save fault")
		}
		if err != nil {
			return err
		}
		// 红包金额支付
		// 1.需要红包中间商的红包资金账户，定义在配置文件中，事先初始化到资金账户表中
		// 2.从红包发送人的资金账户中扣减红包金额
		// 3.将扣减的红包总金额转入红包中间商的红包资金账户

		// 把资金从红包发送人的资金账户里扣除
		body := accountServices.TradeParticipator{
			AccountNo: dto.AccountNo,
			UserId:    dto.UserId,
			Username:  dto.Username,
		}
		systemAccount := base.GetSystemAccount()
		target := accountServices.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			UserId:    systemAccount.UserId,
			Username:  systemAccount.Username,
		}
		transfer := &accountServices.AccountTransferDTO{
			TradeNo:     domain.EnvelopeNo,
			TradeBody:   body,
			TradeTarget: target,
			Amount:      domain.Amount,
			ChangeType:  accountServices.EnvelopeOutgoing,
			ChangeFlag:  accountServices.FlagTransferOut,
			Desc:        "红包金额支付",
		}
		accountDomain := accounts.NewAccountDomain()
		status, err := accountDomain.TransferWithContextTx(ctx, transfer)
		if status != accountServices.TransferredStatusSuccess {
			return err
		}

		// 将扣减的红包总金额转入红包中间商的红包资金账户
		transfer = &accountServices.AccountTransferDTO{
			TradeNo:     domain.EnvelopeNo,
			TradeBody:   target,
			TradeTarget: body,
			Amount:      domain.Amount,
			ChangeType:  accountServices.EnvelopeIncoming,
			ChangeFlag:  accountServices.FlagTransferIn,
			Desc:        "红包金额转入",
		}
		status, err = accountDomain.TransferWithContextTx(ctx, transfer)
		if status != accountServices.TransferredStatusSuccess {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 扣减金额没有问题，返回活动
	activity.RedEnvelopeGoodsDTO = *domain.RedEnvelopeGoods.ToDTO()

	return activity, err
}
