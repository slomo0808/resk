package textx

import (
	"github.com/slomo0808/infra"
	"github.com/slomo0808/infra/base"
	"github.com/slomo0808/infra/comm"
	"github.com/tietang/props/ini"
	"strings"
)

func init() {
	// 获取程序运行文件所在的路径
	path := comm.GetCurrentPath()
	path = strings.TrimRight(path, "textx")
	// 加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(path + "brun/config.ini")

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})

	app := infra.New(conf)
	app.Start()
}
