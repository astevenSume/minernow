// @APIVersion 1.0.0
// @Title logger
// @Description 服务基类定义，用于定义与业务无关的服务基础属性
// @Contact tianguimao@treehousefuture.com
package common

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"os"
	"runtime"
)

var uclogger = NewUCLogger()

type UCLogger struct {
	disableSend2Kafka, isSend2Kafka bool
	retryIntervalSecs               int64
}

func NewUCLogger() *UCLogger {
	return &UCLogger{}
}

func LogInit() error {
	return uclogger.init()
}

func (this *UCLogger) init() (err error) {
	//不管日志目录有没有，尝试创建下
	path, configs := beego.AppConfig.String("log::path"), beego.AppConfig.String("log::configs")
	stdout, err := beego.AppConfig.Bool("log::stdout")
	if err != nil {
		panic(err)
	}

	_ = os.MkdirAll(path, os.ModePerm)
	//日志文件全路径名称
	fileName := fmt.Sprintf("%s/%s.log", path, beego.BConfig.AppName)
	//设置日志文件配置信息。按照日志等级拆分
	if err = beego.SetLogger(logs.AdapterFile, fmt.Sprintf(configs, fileName)); err != nil {
		panic(err)
	}

	//如果禁止控制台输出，则删除默认的stdout日志记录器
	if !stdout {
		if err = beego.BeeLogger.DelLogger(logs.AdapterConsole); err != nil {
			panic(err)
		}
	}

	level, err := beego.AppConfig.Int("log::level")
	if err != nil {
		// default log level set to LevelError
		beego.SetLevel(beego.LevelError)
		err = nil
	} else {
		beego.SetLevel(level)
	}

	return
}

var CurFuncName = func(level int) string {
	if funcName, _, _, ok := runtime.Caller(level); ok {
		return runtime.FuncForPC(funcName).Name()
	}
	return ""
}

var CurFuncLine = func(level int) string {
	if _, file, line, ok := runtime.Caller(level); ok {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return ""
}

//严重错误日志
func LogFuncCritical(format string, args ...interface{}) {
	beego.Critical(formatWithFileLine(format, args...))
}

//错误日志
func LogFuncError(format string, args ...interface{}) {
	beego.Error(formatWithFileLine(format, args...))
}

//告警日志
func LogFuncWarning(format string, args ...interface{}) {
	beego.Warning(formatWithFileLine(format, args...))
}

//跟踪日志
func LogFuncTrace(format string, args ...interface{}) {
	beego.Trace(formatWithFileLine(format, args...))
}

//信息日志
func LogFuncInfo(format string, args ...interface{}) {
	beego.Info(formatWithFileLine(format, args...))
}

//调试日志
func LogFuncDebug(format string, args ...interface{}) {
	beego.Debug(formatWithFileLine(format, args...))
}

func formatWithFileLine(format string, args ...interface{}) string {
	fileLine := CurFuncLine(3)
	message := fmt.Sprintf(format, args...)
	return "[" + fileLine + "] " + message
}
