package main

import (
	"github.com/go-ini/ini"
	_ "imooc.com/resk"
	"imooc.com/resk/comm"
	"imooc.com/resk/infra"
	"log"
)

func main() {
	// 获取程序运行文件所在的路径
	path := comm.GetCurrentPath()
	// 加载和解析配置文件
	conf, err := ini.Load(path + "/config.ini")
	if err != nil {
		log.Fatal(err)
	}
	app := infra.New(conf)
	app.Start()
}
