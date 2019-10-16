package main

import (
	"github.com/tietang/go-eureka-client/eureka"
	"github.com/tietang/props/ini"
	"imooc.com/resk/comm"
	"path/filepath"
)

func main() {
	path := comm.GetCurrentPath()
	conf := ini.NewIniFileConfigSource(filepath.Join(path, "eureka.ini"))
	client := eureka.NewClient(conf)
	client.Start()
	c := make(chan struct{})
	<-c
}
