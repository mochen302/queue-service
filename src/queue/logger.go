package queue

import (
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"sync"
)

var logger *logrus.Logger
var lock sync.Mutex

func LoggerInit(logFilePath string, logFileName string) {
	lock.Lock()
	defer lock.Unlock()

	if logger != nil {
		return
	}

	logger = createLogger(logFilePath, logFileName)
}

func Logger() *logrus.Logger {
	checkLoggerInit()
	return logger
}

func Debug(args ...interface{}) {
	checkLoggerInit()
	logger.Debug(args)
}

func Info(args ...interface{}) {
	checkLoggerInit()
	logger.Info(args)
}

func Warn(args ...interface{}) {
	checkLoggerInit()
	logger.Warn(args)
}

func Error(args ...interface{}) {
	checkLoggerInit()
	logger.Error(args...)
}

func checkLoggerInit() {
	if logger == nil {
		panic("logger has not init!")
	}
}

func createLogger(logFilePath string, logFileName string) *logrus.Logger {
	//日志文件
	fileName := path.Join(logFilePath, logFileName)

	//写入文件
	src, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic("open fileName:" + fileName + " error:" + err.Error())
	}

	//实例化
	logger := logrus.New()

	//设置输出
	logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: false})
	return logger
}
