package main

import (
	"fmt"
	"imooc.com/resk/infra/algo"
)

func main() {
	var count, amount, sum int64 = 10, 10000, 0
	remain := amount
	for i := int64(0); i < count; i++ {
		x := algo.DoubleAverage(count-i, remain)
		remain -= x
		sum += x
		fmt.Printf("%d  ", x)
	}
	fmt.Println()
	fmt.Println(sum)
}
