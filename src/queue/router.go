package queue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Router(r *gin.Engine, q *Queue) {
	/*本不应该用GET方式，只是方便浏览器访问*/
	r.GET("/queue/join", func(c *gin.Context) {

		strId := c.Query("id")
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			setResponse(c, http.StatusBadRequest, fmt.Sprintf("id:%v is not int64", strId)+err.Error())
			return
		}

		nickName := c.Query("nickname")

		handleInternal(c, q.TryJoin, id, nickName)
	})

	/*本不应该用GET方式，只是方便浏览器访问*/
	r.GET("/queue/query", func(c *gin.Context) {

		strId := c.Query("id")
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			setResponse(c, http.StatusBadRequest, fmt.Sprintf("id:%v is not int64", strId)+err.Error())
			return
		}

		handleInternal(c, q.QueryState, id)
	})

	r.GET("/queue/stat", func(c *gin.Context) {
		handleInternal(c, q.StatInfo)
	})

}

func handleInternal(c *gin.Context, f func(p ...interface{}) (result interface{}), param ...interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			setResponse(c, http.StatusInternalServerError, fmt.Sprintf(err.(string)))
		}
	}()

	result := f(param)
	setResponse(c, http.StatusOK, "success", fmt.Sprint(result))
}

func setResponse(c *gin.Context, code int, message string, result ...string) {
	if len(result) > 0 {
		c.JSON(code, gin.H{
			"message": message,
			"result":  result[0],
		})
	} else {
		c.JSON(code, gin.H{
			"message": message,
		})
	}
}
