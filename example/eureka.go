package main

import (
	"github.com/tietang/go-eureka-client/eureka"
	"time"
)

func main() {
	cfg := eureka.Config{
		CertFile:    "",
		KeyFile:     "",
		CaCertFile:  nil,
		DialTimeout: time.Second * 10,
		Consistency: "",
	}
	client := eureka.NewClientByConfig([]string{
		"http://127.0.0.1:8761/eureka",
	}, cfg)
	appName := "Go-Example"
	instance := eureka.NewInstanceInfo("test.com",
		appName,
		"127.0.0.2",
		8080, 30,
		false)
	client.RegisterInstance(appName, instance)
	client.Start() //该方法不阻塞
	waitChan := make(chan struct{})
	<-waitChan
}
