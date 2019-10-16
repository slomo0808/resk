package web

import (
	"github.com/kataras/iris"
	"github.com/slomo0808/infra"
	"github.com/slomo0808/infra/base"
	"imooc.com/resk/services"
)

func init() {
	infra.RegisterApi(&AccountApi{})
}

// 定义web api的时候，对每一个子业务，定义统一的前缀
// 资金账户的根路径定义为 ：/account
// 版本号： /v1/account
type AccountApi struct{}

func (a *AccountApi) Init() {
	groupRouter := base.Iris().Party("/v1/account")
	groupRouter.Post("/create", createHandler)
	groupRouter.Post("/transfer", transferHandler)
	groupRouter.Get("/envelope/get", getEnvelopeAccountHandler)
	groupRouter.Get("/get", getAccountHandler)
}

// 账户创建接口 /v1/account/create
func createHandler(ctx iris.Context) {
	// 获取请求参数
	account := &services.AccountCreatedDTO{}
	err := ctx.ReadJSON(account)
	r := &base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestPaamsErr
		r.Message = err.Error()
		ctx.JSON(r)
		return
	}
	// 执行创建账户的代码
	service := services.GetAccountService()
	dto, err := service.CreateAccount(account)
	if err != nil {
		r.Code = base.ResCodeIntenalServerErr
		r.Message = err.Error()
	}
	r.Data = dto
	ctx.JSON(r)
}

// 转账接口 /v1/account/transfer
func transferHandler(ctx iris.Context) {
	// 获取请求参数
	account := &services.AccountTransferDTO{}
	err := ctx.ReadJSON(account)
	r := &base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestPaamsErr
		r.Message = err.Error()
		ctx.JSON(r)
		return
	}
	// 执行转账逻辑
	service := services.GetAccountService()
	status, err := service.Transfer(account)
	if err != nil {
		r.Code = base.ResCodeIntenalServerErr
		r.Message = err.Error()
	}
	r.Data = status
	if status != services.TransferredStatusSuccess {
		r.Code = base.ResCodeBizTransferredFailure
		r.Message = err.Error()
	}
	ctx.JSON(r)
}

// 查询红包账户的接口 v1/account/envelope/get
func getEnvelopeAccountHandler(ctx iris.Context) {
	userId := ctx.URLParam("user_id")
	service := services.GetAccountService()
	dto := service.GetEnvelopeAccountByUserId(userId)
	r := &base.Res{
		Code: base.ResCodeOk,
	}
	if dto == nil {
		r.Code = base.ResCodeValidationErr
		r.Message = "没有查询到数据"
		ctx.JSON(r)
		return
	}
	r.Data = dto
	ctx.JSON(r)
}

// 查询账户信息的接口 v1/account/get
func getAccountHandler(ctx iris.Context) {
	accountNo := ctx.URLParam("account_no")
	service := services.GetAccountService()
	dto := service.GetAccount(accountNo)
	r := &base.Res{
		Code: base.ResCodeOk,
	}
	if dto == nil {
		r.Code = base.ResCodeValidationErr
		r.Message = "没有查询到数据"
		ctx.JSON(r)
		return
	}
	r.Data = dto
	ctx.JSON(r)
}
