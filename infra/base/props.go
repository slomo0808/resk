package base

import (
	"fmt"
	"github.com/go-ini/ini"
	"imooc.com/resk/infra"
	"sync"
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
	GetSystemAccount()
	fmt.Println("初始化配置。")
}

type SystemAccount struct {
	AccountNo   string
	AccountName string
	UserId      string
	Username    string
}

var systemAccount *SystemAccount
var systemAccountOnce sync.Once

func GetSystemAccount() *SystemAccount {
	systemAccountOnce.Do(func() {
		systemAccount = &SystemAccount{
			AccountNo:   Props().Section("system.account").Key("accountNo").MustString("10000020190101010000000000000001"),
			AccountName: Props().Section("system.account").Key("accountName").MustString("系统红包账户"),
			UserId:      Props().Section("system.account").Key("userId").MustString("10001"),
			Username:    Props().Section("system.account").Key("username").MustString("系统红包账户"),
		}
	})
	return systemAccount
}

func GetEnvelopeActivityLink() string {
	return Props().Section("envelope").Key("link").MustString("/v1/envelope/link")
}

func GetEnvelopeDomain() string {
	return Props().Section("envelope").Key("domain").MustString("http://localhost")
}
