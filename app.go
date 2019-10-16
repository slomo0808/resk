package resk

import (
	_ "github.com/slomo0808/account"
	"github.com/slomo0808/infra"
	"github.com/slomo0808/infra/base"
	"imooc.com/resk/apis/gorpc"
	_ "imooc.com/resk/apis/web"
	_ "imooc.com/resk/core/envelopes"
	"imooc.com/resk/jobs"
	_ "imooc.com/resk/views"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.GoRPCStarter{})
	infra.Register(&gorpc.GoRPCApiStarter{})
	infra.Register(&jobs.RefundExpiredJobStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
	infra.Register(&base.EurekaStarter{})
	infra.Register(&base.HookStarter{})
}
