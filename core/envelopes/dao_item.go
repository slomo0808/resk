package envelopes

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type RedEnvelopeItemDao struct {
	runner *dbx.TxRunner
}

func (dao *RedEnvelopeItemDao) GetOne(itemNo int64) *RedEnvelopeItem {
	out := &RedEnvelopeItem{ItemNo: itemNo}
	ok, err := dao.runner.GetOne(out)
	if !ok || err != nil {
		logrus.Error(err)
		return nil
	}
	return out
}

func (dao *RedEnvelopeItemDao) Insert(data *RedEnvelopeItem) (int64, error) {
	res, err := dao.runner.Insert(data)
	if err != nil {
		return 0, nil
	}
	return res.LastInsertId()
}
