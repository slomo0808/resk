package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"imooc.com/resk/services"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", ":18082")
	if err != nil {
		log.Panic(err)
	}
	in := services.RedEnvelopeSendingDTO{
		EnvelopeType: services.LuckyEnvelopeType,
		Username:     "rpc测试用户",
		UserId:       "1P3z3B5bMhKB45aBSJizAXYWgkS",
		Blessing:     "rpc测试",
		Amount:       decimal.NewFromFloat(10),
		Quantity:     10,
	}
	out := &services.RedEnvelopeActivity{}
	err = client.Call("EnvelopeRpc.SendOut", in, out)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v", out)
}
