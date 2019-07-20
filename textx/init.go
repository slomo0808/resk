package textx

import (
	"github.com/go-ini/ini"
	"imooc.com/resk/comm"
	"imooc.com/resk/infra"
	"imooc.com/resk/infra/base"
	"log"
	"strings"
)

func init() {
	// 获取程序运行文件所在的路径
	path := comm.GetCurrentPath()
	path = strings.TrimRight(path, "textx")
	// 加载和解析配置文件
	conf, err := ini.Load(path + "brun/config.ini")
	if err != nil {
		log.Fatal(err)
	}

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})

	app := infra.New(conf)
	app.Start()
}
