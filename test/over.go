package test

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
	"github.com/tietang/dbx"
	"log"
)

var db *dbx.Database

func init() {
	settings := dbx.Settings{
		DriverName: "mysql",
		User:       "root",
		Password:   "123456",
		Database:   "po0",
		Host:       "127.0.0.1:3306",
		Options: map[string]string{
			"parseTime": "true",
		},
	}
	var err error
	if db, err = dbx.Open(settings); err != nil {
		log.Fatal("over.dbx.Open error:", err)
	}
	db.RegisterTable(&GoodsSigned{}, "goods")
	db.RegisterTable(&GoodsUnsigned{}, "goods_unsigned")
}

// 事务所方案
func UpdateForLock(g Goods) {
	// 通过db.Tx 函数构建事务所代码块
	err := db.Tx(func(runner *dbx.TxRunner) error {
		// 第一步：锁定需要修改的资源，也就是需要修改的数据行
		// 编写事务所查询语句, 使用for update祖居来锁定资源
		query := "select * from goods where envelope_no = ? for update"
		out := &GoodsSigned{}
		ok, err := runner.Get(out, query, g.EnvelopeNo)
		if !ok || err != nil {
			return fmt.Errorf("UpdateForLock.db.Tx.runner.Get ok = %t, error = %s", ok, err)
		}
		// 第二步：计算剩余金额和剩余数量
		subAmount := decimal.NewFromFloat(0.01)
		remainAmount := out.RemainAmount.Sub(subAmount)
		remainQuantity := out.RemainQuantity - 1
		// 第三步：执行更新
		update := "update goods set remain_amount = ?, remain_quantity=? " +
			" where envelope_no = ?"
		_, rows, err := db.Execute(update, remainAmount, remainQuantity, g.EnvelopeNo)
		if err != nil || rows < 1 {
			return fmt.Errorf("UpdateForLock.db.Tx.db.Execute rows = %d, error = %s", rows, err)
		}
		return nil
	})
	if err != nil {
		fmt.Println("UpdateForLock.db.Tx error = ", err)
	}
}

// 数据库无符号类型+直接更新方案
func UpdateForUnsigned(g Goods) {
	update := "update goods_unsigned set remain_amount=remain_amount-?," +
		"remain_quantity=remain_quantity-? where envelope_no=?"
	_, rows, err := db.Execute(update, 0.01, 1, g.EnvelopeNo)
	if err != nil {
		fmt.Println(err)
	}
	if rows < 1 {
		fmt.Println("扣减失败")
	}
}

// 乐观锁方案
func UpdateForOptimistic(g Goods) {
	update := "update goods set remain_amount=remain_amount-?," +
		" remain_quantity=remain_quantity-? where envelope_no=? " +
		" and remain_amount>=? " +
		" and remain_quantity>=? "
	_, rows, err := db.Execute(update, 0.01, 1, g.EnvelopeNo, 0.01, 1)
	if err != nil {
		fmt.Println(err)
	}
	if rows < 1 {
		fmt.Println("扣减失败")
	}
}

// 乐观锁方案 and 无符号字段 双保险
func UpdateForOptimisticAndUnsigned(g Goods) {
	update := "update goods_unsigned set remain_amount=remain_amount-?," +
		" remain_quantity=remain_quantity-? where envelope_no=? " +
		" and remain_amount>=? " +
		" and remain_quantity>=? "
	_, rows, err := db.Execute(update, 0.01, 1, g.EnvelopeNo, 0.01, 1)
	if err != nil {
		fmt.Println(err)
	}
	if rows < 1 {
		fmt.Println("扣减失败")
	}
}
