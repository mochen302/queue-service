package main

import (
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

const (
	LOG_PATH            = "./output/"
	LOG_FILE_NAME       = "server.log"
	QUEUE_MAX_USER      = 100 * 10000
	QUEUE_MAX_USER_WAIT = 20 * 1000
)

func main() {
	testTreeMap()

	logrus1 := createLogger(LOG_PATH, LOG_FILE_NAME)

	r := gin.Default()
	r.Use(LoggerToFile(logrus1))

	r.GET("/get/:key/:value", func(c *gin.Context) {

		key := c.Param("key")
		value := c.Param("value")
		logrus1.Info(key, "->", value)

		c.JSON(200, gin.H{
			"message": "pong",
			key:       value,
		})
	})

	err := r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		panic("start server at:localhost:8080 error" + err.Error())
	}
}

func testTreeMap() {
	m := treemap.NewWithIntComparator() // empty (keys are of type int)
	m.Put(1, "x")                       // 1->x
	m.Put(2, "b")                       // 1->x, 2->b (in order)
	m.Put(1, "a")                       // 1->a, 2->b (in order)
	_, _ = m.Get(2)                     // b, true
	_, _ = m.Get(3)                     // nil, false
	_ = m.Values()                      // []interface {}{"a", "b"} (in order)
	_ = m.Keys()                        // []interface {}{1, 2} (in order)
	m.Remove(1)                         // 2->b
	m.Clear()                           // empty
	m.Empty()                           // true
	m.Size()                            // 0cd

	// Other:
	m.Min() // Returns the minimum key and its value from map.
	m.Max() // Returns the maximum key and its value from map.
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
