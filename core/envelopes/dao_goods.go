package envelopes

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"imooc.com/resk/services"
	"time"
)

type RedEnvelopeDao struct {
	runner *dbx.TxRunner
}

// 插入
func (dao *RedEnvelopeDao) Insert(po *RedEnvelopeGoods) (int64, error) {
	rs, err := dao.runner.Insert(po)
	if err != nil {
		logrus.Error(err)
		return 0, nil
	}
	return rs.LastInsertId()
}

// 查询，根据红包编号
func (dao *RedEnvelopeDao) GetOne(envelopeNo string) *RedEnvelopeGoods {
	var out = &RedEnvelopeGoods{EnvelopeNo: envelopeNo}
	ok, err := dao.runner.GetOne(out)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return out
}

// 更新红包库存
// 不再使用事务行锁来更新红包余额和数量
// 改用乐观锁来保证更新操作的安全性，避免负库存问题
// 通过where子句来判断红包剩余金额和数量来解决两个问题
// 1. 负库存问题，避免红包金额和数量不足时任然进行扣减
// 2. 减少实际的库存更新，也就是过滤掉无效的更新，提高总体性能
func (dao *RedEnvelopeDao) UpdateBalance(envelopeNo string, amount decimal.Decimal) (int64, error) {
	sqlQuery := "update red_envelope_goods " +
		" set remain_amount = remain_amount - CAST(? as DECIMAL(30,6)), " +
		" remain_quantity = remain_quantity - 1 " +
		" where envelope_no = ? " +
		" and remain_quantity > 0 " +
		" and remain_amount >= CAST(? as DECIMAL(30,6)) "
	res, err := dao.runner.Exec(sqlQuery, amount.String(), envelopeNo, amount.String())
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return res.RowsAffected()
}

// 退款更新订单状态
func (dao *RedEnvelopeDao) UpdateOrderStatus(envelopeNo string, status services.OrderStatus) (int64, error) {
	sqlQuery := "update red_envelope_goods " +
		" set status = ? " +
		" where envelope_no = ?"
	res, err := dao.runner.Exec(sqlQuery, status, envelopeNo)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return res.RowsAffected()
}

// 过期，把所有的过期红包都查询出来，分页，limit offset size
func (dao *RedEnvelopeDao) FindExpired(offset, size int) []RedEnvelopeGoods {
	var goods = make([]RedEnvelopeGoods, 0)
	now := time.Now()
	sqlQuery := "select * from red_envelope_goods " +
		" where expired_at < ? " +
		" limit ?,?"
	err := dao.runner.Find(&goods, sqlQuery, now, offset, size)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}
