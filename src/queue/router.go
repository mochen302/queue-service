package queue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type BUSINESS_CODE int8

const (
	SUCCESS      BUSINESS_CODE = 0
	EXCEPTION    BUSINESS_CODE = -1
	ID_NOT_VALID BUSINESS_CODE = -2
	ID_LESS_ZERO BUSINESS_CODE = -3
)

func Router(r *gin.Engine, q *Queue) {
	/*本不应该用GET方式，只是方便浏览器访问*/
	r.GET("/queue/join", func(c *gin.Context) {

		strId := c.Query("id")
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			setResponse(c, http.StatusOK, fmt.Sprintf("id:%v is not int64", strId)+err.Error(), ID_NOT_VALID)
			return
		}

		if id <= 0 {
			setResponse(c, http.StatusOK, fmt.Sprintf("id:%v is <=0 ", strId), ID_LESS_ZERO)
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
			setResponse(c, http.StatusOK, fmt.Sprintf("id:%v is not int64", strId)+err.Error(), ID_NOT_VALID)
			return
		}

		if id <= 0 {
			setResponse(c, http.StatusOK, fmt.Sprintf("id:%v is <=0 ", strId), ID_LESS_ZERO)
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
			Error(err)
			setResponse(c, http.StatusInternalServerError, "internal error", EXCEPTION)
		}
	}()

	result := f(param)
	setResponse(c, http.StatusOK, "success", SUCCESS, fmt.Sprint(result))
}

func setResponse(c *gin.Context, code int, message string, businessCode BUSINESS_CODE, result ...string) {
	if len(result) > 0 {
		c.JSON(code, gin.H{
			"code":    businessCode,
			"message": message,
			"result":  result[0],
		})
	} else {
		c.JSON(code, gin.H{
			"code":    businessCode,
			"message": message,
			"result":  "",
		})
	}
}
