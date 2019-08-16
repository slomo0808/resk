package gorpc

import (
	"imooc.com/resk/infra"
	"imooc.com/resk/infra/base"
)

type GoRPCApiStarter struct {
	infra.BaseStarter
}

func (s *GoRPCApiStarter) Init(ctx infra.StarterContext) {
	base.RpcRegister(new(EnvelopeRpc))
}
