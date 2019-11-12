package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mochen302/queue-service/src/queue"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	LOG_PATH         = "./output/"
	LOG_FILE_NAME    = "server.log"
	ADDRESS          = "127.0.0.1:8080"
	MAX_HANDLE_COUNT = 100
	MAX_WAIT_COUNT   = 20 * 10000
)

func main() {

	queue.LoggerInit(LOG_PATH, LOG_FILE_NAME)

	queueService := queue.New(MAX_HANDLE_COUNT, MAX_WAIT_COUNT)

	r := gin.Default()
	r.Use(LoggerToFile(queue.Logger()))
	queue.Router(r, queueService)
	err := r.Run(ADDRESS)
	if err != nil {
		panic("start server at:" + ADDRESS + " error" + err.Error())
	}
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
