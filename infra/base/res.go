package base

type ResCode int

const (
	ResCodeOk                    = 1000
	ResCodeValidationErr         = 2000
	ResCodeRequestPaamsErr       = 2100
	ResCodeIntenalServerErr      = 5000
	ResCodeBizErr                = 6000
	ResCodeBizTransferredFailure = 6010
)

type Res struct {
	Code    ResCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
