package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func main() {

	logFilePath := "/Users/wuqiong/mochen/gopath/src/github.com/mochen302/queue-service/output/"
	logFileName := "server.log"
	logrus := createLogger(logFilePath, logFileName)

	r := gin.Default()
	r.Use(LoggerToFile(logrus))

	r.GET("/get/:key/:value", func(c *gin.Context) {

		key := c.Param("key")
		value := c.Param("value")
		logrus.Info(key, "->"+value)

		c.JSON(200, gin.H{
			"message": "pong",
			key:       value,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func createLogger(logFilePath string, logFileName string) *logrus.Logger {
	//日志文件
	fileName := path.Join(logFilePath, logFileName)

	//写入文件
	src, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//实例化
	logger := logrus.New()

	//设置输出
	logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{})
	return logger
}

func LoggerToFile(logger *logrus.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
