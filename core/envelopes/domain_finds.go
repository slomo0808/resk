package envelopes

import (
	"github.com/tietang/dbx"
	"imooc.com/resk/infra/base"
)

func (domain *goodsDomain) Find(po *RedEnvelopeGoods, offset, limit int) (regs []RedEnvelopeGoods) {
	base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeDao{runner: runner}
		regs = dao.Find(po, offset, limit)
		return nil
	})
	return regs
}

func (domain *goodsDomain) FindByUser(userId string, offset, limit int) (regs []RedEnvelopeGoods) {
	base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeDao{runner: runner}
		regs = dao.FindByUser(userId, offset, limit)
		return nil
	})
	return regs
}
func (domain *goodsDomain) GetOne(envelopeNo string) (po *RedEnvelopeGoods) {
	base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeDao{runner: runner}
		po = dao.GetOne(envelopeNo)
		return nil
	})
	return po
}

// 查询可以收的红包
func (domain *goodsDomain) ListReceivable(offset, limit int) (regs []RedEnvelopeGoods) {
	base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeDao{runner: runner}
		regs = dao.ListReceivable(offset, limit)
		return nil
	})
	return regs
}

// 查询已收到的红包
func (domain *goodsDomain) ListReceived(userId string, offset, limit int) (regs []*RedEnvelopeItem) {
	base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		regs = dao.ListReceivedItems(userId, offset, limit)
		return nil
	})
	return regs
}
