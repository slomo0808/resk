package accounts

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountDao struct {
	runner *dbx.TxRunner
}

// 查询数据库持久化对象单实例， 即获取一行数据
func (dao *AccountDao) GetOne(accountNo string) *Account {
	out := &Account{AccountNo: accountNo}
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

// 通过用户Id和账户类型来查询账户信息
func (dao *AccountDao) GetByUserId(userId string, accountType int) *Account {
	out := &Account{}
	sqlQuery := "select * from account where user_id=? and account_type=?"
	ok, err := dao.runner.Get(out, sqlQuery, userId, accountType)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return out
}

// 账户数据的插入
func (dao *AccountDao) Insert(data *Account) (int64, error) {
	result, err := dao.runner.Insert(data)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	return result.LastInsertId()
}

// 账户余额的更新
// amount > 0 收入 amount < 0 扣减
func (dao *AccountDao) UpdateBalance(accountNo string, amount decimal.Decimal) (rows int64, err error) {
	sqlQuery := "update account " +
		" set balance=balance+CAST(? as DECIMAL(30,6)) " +
		" where account_no=? " +
		" and balance>=-1*CAST(? as DECIMAL(30,6))"
	rs, err := dao.runner.Exec(sqlQuery, amount.String(), accountNo, amount.String())
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

// 账户状态的更新
func (dao *AccountDao) UpdateStatus(accountNo string, status int) (rows int64, err error) {
	sqlQuery := "update account set status = ? where account_no = ?"
	rs, err := dao.runner.Exec(sqlQuery, status, accountNo)
	if err != nil {
		logrus.Error(err)
		return 0, nil
	}
	return rs.RowsAffected()
}
