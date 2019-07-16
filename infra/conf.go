package infra

import "fmt"

type ConfStarter struct {
	BaseStarter
}

func init() {
	Register(&ConfStarter{})
}

func (c *ConfStarter) Init(ctx StarterContext) {
	fmt.Println("配置初始化")
}

func (c *ConfStarter) Setup(ctx StarterContext) {
	fmt.Println("配置安装")
}

func (c *ConfStarter) Start(ctx StarterContext) {
	fmt.Println("配置启动")
}
