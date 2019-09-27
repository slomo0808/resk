package base

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/tietang/props/kvs"
	"imooc.com/resk/infra"
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
	logrus.Info("DbxDatabaseStarter Setup()")
	conf := ctx.Props()
	// 数据库配置
	settings := dbx.Settings{}
	err := kvs.Unmarshal(conf, &settings, "mysql")
	if err != nil {
		panic(err)
	}
	logrus.Infof("%+v", settings)
	logrus.Info("mysql.conn url:", settings.ShortDataSourceName())

	dbconn, err := dbx.Open(settings)
	if err != nil {
		logrus.Panic("dbx.Setup dbx.Open error:", err)
	}
	logrus.Info(dbconn.Ping())
	database = dbconn
}
