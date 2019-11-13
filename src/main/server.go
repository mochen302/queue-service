package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mochen302/queue-service/src/queue"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const (
	LOG_PATH         = "./output/"
	LOG_FILE_NAME    = "server.log"
	ADDRESS          = ":8080"
	MAX_HANDLE_COUNT = 100
	MAX_WAIT_COUNT   = 20 * 10000
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	queue.LoggerInit(LOG_PATH, LOG_FILE_NAME)

	queueService := queue.New(MAX_HANDLE_COUNT, MAX_WAIT_COUNT)

	gin := gin.Default()
	gin.Use(LoggerToFile(queue.Logger()))

	queue.Router(gin, queueService)

	err := gin.Run(ADDRESS)
	if err != nil {
		panic("start server at:" + ADDRESS + " error" + err.Error())
	}

	server := &http.Server{
		Addr:         ADDRESS,
		Handler:      gin,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	gracefulExitWeb(queueService, server)

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

func gracefulExitWeb(queueService *queue.Queue, server *http.Server) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-ch
	queue.Error("got a signal", sig)

	now := time.Now()

	queueService.Close()

	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(cxt)
	if err != nil {
		queue.Error("err", err)
	}

	queue.Error("shutdown server cost:", time.Since(now))
}
