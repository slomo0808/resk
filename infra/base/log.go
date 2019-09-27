package base

import (
	"github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	// 定义日志格式
	formatter := &prefixed.TextFormatter{}
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05.000000"
	formatter.ForceFormatting = true
	formatter.ForceColors = true
	formatter.DisableColors = false
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red",
		PanicLevelStyle: "red",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "black+h",
	})
	log.SetFormatter(formatter)
	// 日志级别
	level := os.Getenv("log.debug")
	if level == "true" {
		log.SetLevel(log.DebugLevel)
	}
	log.SetLevel(log.DebugLevel)
	// 控制台高亮显示
	//log.Info("测试")
	//log.Debug("debug")

	// 显示文件名代码行数
	log.SetReportCaller(true)

	// 日志文件和滚动配置
	// github.com/lestrrat/go-file-rotatelogs
	logFileSettings()
}

func logFileSettings() {
	// 配置日志输入目录
	logPath, _ := filepath.Abs("./logs")
	log.Infof("log dir: %s", logPath)
	logFileName := "resk"
	// 日志文件最大保存时间 24h
	maxAge := 24 * time.Hour
	// 日志切割时间间隔 1h
	rotationTime := time.Hour * 1

	os.MkdirAll(logPath, os.ModePerm)

	baseLogPath := path.Join(logPath, logFileName)
	// 设置滚动日志输出
	writer, err := rotatelogs.New(
		strings.TrimSuffix(baseLogPath, ".log")+".%Y%m%d%H.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime)) // 日志切割时间间隔

	if err != nil {
		log.Errorf("config local file system logger error = %+v", err)
	}

	writers := io.MultiWriter(writer, os.Stdout)
	log.SetOutput(writers)
}
