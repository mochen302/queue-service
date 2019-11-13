package queue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Router(r *gin.Engine, queueService *Queue) {
	r.GET("/queue/join", func(c *gin.Context) {

		strId := c.Query("id")
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			setResponse(c, http.StatusBadRequest, fmt.Sprintf("id:%v is not int64", strId)+err.Error())
			return
		}

		nickName := c.Query("nickname")

		func() {
			defer func() {
				err := recover()
				if err != nil {
					setResponse(c, http.StatusInternalServerError, fmt.Sprintf(err.(string)))
				}
			}()

			result := queueService.TryJoin(id, nickName)
			setResponse(c, http.StatusOK, "success", fmt.Sprint(result))
		}()

	})

	r.GET("/queue/query", func(c *gin.Context) {

		strId := c.Query("id")
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			setResponse(c, http.StatusBadRequest, fmt.Sprintf("id:%v is not int64", strId)+err.Error())
			return
		}

		func() {
			defer func() {
				err := recover()
				if err != nil {
					setResponse(c, http.StatusInternalServerError, fmt.Sprintf(err.(string)))
				}
			}()

			result := queueService.QueryState(id)
			setResponse(c, http.StatusOK, "success", fmt.Sprint(result))
		}()

	})

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
