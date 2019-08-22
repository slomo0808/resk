package envelopes

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type RedEnvelopeItemDao struct {
	runner *dbx.TxRunner
}

func (dao *RedEnvelopeItemDao) GetOne(itemNo string) *RedEnvelopeItem {
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

func (dao *RedEnvelopeItemDao) FindItems(envelopeNo string) []*RedEnvelopeItem {
	items := make([]*RedEnvelopeItem, 0)
	sqlQuery := "select * from red_envelope_item where envelope_no = ?"
	err := dao.runner.Find(items, sqlQuery, envelopeNo)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return items
}
