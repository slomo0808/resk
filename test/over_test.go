package test

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"log"
	"testing"
)

// 基准测试代码，Benchmark开头，第一个字母B大写
// testing.B

// 事务锁
func BenchmarkUpdateForLock(b *testing.B) {
	g := GoodsSigned{}
	g.EnvelopeNo = ksuid.New().Next().String()
	g.RemainQuantity = 100000
	g.RemainAmount = decimal.NewFromFloat(10000000)
	_, err := db.Insert(g)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		UpdateForLock(g.Goods)
	}
}

// 乐观锁+无符号双保险
func BenchmarkUpdateForOptimisticAndUnsigned(b *testing.B) {
	g := GoodsUnsigned{}
	g.EnvelopeNo = ksuid.New().Next().String()
	g.RemainQuantity = 100000
	g.RemainAmount = decimal.NewFromFloat(10000000)
	_, err := db.Insert(g)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		UpdateForOptimisticAndUnsigned(g.Goods)
	}
}

// 基准测试：无符号字段类型+直接更新
func BenchmarkUpdateForUnsigned(b *testing.B) {
	g := GoodsUnsigned{}
	g.EnvelopeNo = ksuid.New().Next().String()
	g.RemainQuantity = 100000
	g.RemainAmount = decimal.NewFromFloat(10000000)
	_, err := db.Insert(g)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		UpdateForUnsigned(g.Goods)
	}
}

// 基准测试：乐观锁
func BenchmarkUpdateForOptimistic(b *testing.B) {
	g := GoodsSigned{}
	g.EnvelopeNo = ksuid.New().Next().String()
	g.RemainQuantity = 100000
	g.RemainAmount = decimal.NewFromFloat(10000000)
	_, err := db.Insert(g)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		UpdateForOptimistic(g.Goods)
	}
}
