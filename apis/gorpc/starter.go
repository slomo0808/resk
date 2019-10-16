package gorpc

import (
	"github.com/slomo0808/infra"
	"github.com/slomo0808/infra/base"
)

type GoRPCApiStarter struct {
	infra.BaseStarter
}

func (s *GoRPCApiStarter) Init(ctx infra.StarterContext) {
	base.RpcRegister(new(EnvelopeRpc))
}
