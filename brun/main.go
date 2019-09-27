package main

import (
	"github.com/tietang/props/ini"
	_ "imooc.com/resk"
	"imooc.com/resk/comm"
	"imooc.com/resk/infra"
)

func main() {
	// 获取程序运行文件所在的路径
	path := comm.GetCurrentPath()
	// 加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(path + "/config.ini")

	app := infra.New(conf)
	app.Start()
}
