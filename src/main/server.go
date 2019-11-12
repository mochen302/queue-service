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
	SERVER_PORT      = "8080"
	MAX_HANDLE_COUNT = 100
	MAX_WAIT_COUNT   = 20 * 10000
)

func main() {

	queue.LoggerInit(LOG_PATH, LOG_FILE_NAME)

	r := gin.Default()
	r.Use(LoggerToFile(queue.Logger()))

	r.POST("/get/:key/:value", func(c *gin.Context) {

		key := c.Param("key")
		value := c.Param("value")

		c.JSON(200, gin.H{
			"message": "pong",
			key:       value,
		})
	})

	err := r.Run(SERVER_PORT) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		panic("start server at:localhost:8080 error" + err.Error())
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
