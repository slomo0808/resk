package base

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"imooc.com/resk/infra"
	"log"
	"time"
)

// dbx 数据库实例
var database *dbx.Database

func DbxDatabase() *dbx.Database {
	return database
}

// dbx数据库starter，并且设置为全局
type DbxDatabaseStarter struct {
	infra.BaseStarter
}

func (s *DbxDatabaseStarter) Setup(ctx infra.StarterContext) {
	conf := ctx.Props()
	// 数据库配置
	settings := dbx.Settings{
		DriverName: conf.Section("mysql").Key("driverName").String(),
		User:       conf.Section("mysql").Key("username").String(),
		Password:   conf.Section("mysql").Key("password").String(),
		Database:   conf.Section("mysql").Key("database").String(),
		Host:       conf.Section("mysql").Key("host").String(),
		Options: map[string]string{
			"charset":   "utf8",
			"parseTime": "true",
			"loc":       "Local",
		},
		ConnMaxLifetime: conf.Section("mysql").Key("connMaxLifetime").MustDuration(7 * time.Hour),
		MaxOpenConns:    conf.Section("mysql").Key("maxOpenConns").MustInt(5),
		MaxIdleConns:    conf.Section("mysql").Key("maxIdleConns").MustInt(2),
	}
	logrus.Infof("%+v", settings)
	logrus.Info("mysql.conn url:", settings.ShortDataSourceName())

	dbconn, err := dbx.Open(settings)
	if err != nil {
		log.Fatal("dbx.Setup dbx.Open error:", err)
	}
	logrus.Info(dbconn.Ping())
	database = dbconn
}
