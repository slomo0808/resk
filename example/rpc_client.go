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
	send(client)
	receive(client)

}

func send(c *rpc.Client) {
	in := services.RedEnvelopeSendingDTO{
		EnvelopeType: services.LuckyEnvelopeType,
		Username:     "rpc测试用户",
		UserId:       "rpcrpcrpcrpcrpcrpc",
		Blessing:     "rpc测试",
		Amount:       decimal.NewFromFloat(10),
		Quantity:     10,
	}
	out := &services.RedEnvelopeActivity{}
	err := c.Call("EnvelopeRpc.SendOut", in, out)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v\n", out)
}

func receive(c *rpc.Client) {
	in2 := services.RedEnvelopeReceiveDTO{
		EnvelopeNo:   "",
		RecvUsername: "",
		RecvUserId:   "",
	}
	out2 := &services.RedEnvelopeItemDTO{}
	err := c.Call("EnvelopeRpc.Receive", in2, out2)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v\n", out2)
}
