package algo

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDoubleAverage(t *testing.T) {
	var count, amount, sum int64 = 10, 10000, 0
	remain := amount
	for i := int64(0); i < count; i++ {
		x := DoubleAverage(count-i, remain)
		remain -= x
		sum += x
	}
	convey.Convey("二倍均值算法", t, func() {
		convey.So(sum, convey.ShouldEqual, amount)
	})
}
