package base

import (
	"fmt"
	"github.com/go-ini/ini"
	"imooc.com/resk/infra"
)

var props *ini.File

func Props() *ini.File {
	return props
}

type PropsStarter struct {
	infra.BaseStarter
}

func (p *PropsStarter) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	fmt.Println("初始化配置。")
}
