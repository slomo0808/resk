package web

import (
	"github.com/kataras/iris"
	"github.com/slomo0808/infra"
	"github.com/slomo0808/infra/base"
	"imooc.com/resk/services"
)

type RedEnvelopeApi struct {
	service services.RedEnvelopeService
}

func init() {
	infra.RegisterApi(&RedEnvelopeApi{})
}

func (api *RedEnvelopeApi) Init() {
	api.service = services.GetRedEnvelopeService()
	groupParty := base.Iris().Party("/v1/envelope")
	groupParty.Post("/sendout", api.sendOutHandler)
	groupParty.Post("/receive", api.receiveHandler)
}

/*
{
	"envelope_type":0,
	"username":"",
	"user_id":"",
	"amount":"0",
	"quantity":0
}

*/
func (api *RedEnvelopeApi) sendOutHandler(ctx iris.Context) {
	dto := &services.RedEnvelopeSendingDTO{}
	err := ctx.ReadJSON(dto)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestPaamsErr
		r.Message = err.Error()
		ctx.JSON(r)
		return
	}
	// 执行发红包代码
	activity, err := api.service.SendOut(*dto)
	if err != nil {
		r.Code = base.ResCodeRequestPaamsErr
		r.Message = err.Error()
		ctx.JSON(r)
		return
	}
	r.Data = activity
	ctx.JSON(r)
}

func (api *RedEnvelopeApi) receiveHandler(ctx iris.Context) {
	dto := &services.RedEnvelopeReceiveDTO{}
	err := ctx.ReadJSON(dto)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestPaamsErr
		r.Message = err.Error()
		ctx.JSON(r)
		return
	}
	// 执行发红包代码
	item, err := api.service.Receive(dto)
	if err != nil {
		r.Code = base.ResCodeRequestPaamsErr
		r.Message = err.Error()
		ctx.JSON(r)
		return
	}
	r.Data = item
	ctx.JSON(r)
}
