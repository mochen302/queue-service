package queue

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Router(r *gin.Engine, queueService *QueueService) {
	r.POST("/queue/join", func(c *gin.Context) {

		strId := c.Param("id")
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			setResponse(c, 500, fmt.Sprintf("id:%v is not int64", strId))
			return
		}

		nickName := c.Param("nickname")

		func() {
			defer func() {
				err := recover()
				if err != nil {
					setResponse(c, 500, fmt.Sprintf(err.(string)))
				}
			}()

			result := queueService.TryJoin(id, nickName)
			setResponse(c, 200, "success", fmt.Sprint(result))
		}()

	})

	r.POST("/queue/query", func(c *gin.Context) {

		strId := c.Param("id")
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			setResponse(c, 500, fmt.Sprintf("id:%v is not int64", strId))
			return
		}

		func() {
			defer func() {
				err := recover()
				if err != nil {
					setResponse(c, 500, fmt.Sprintf(err.(string)))
				}
			}()

			result := queueService.QueryState(id)
			setResponse(c, 200, "success", fmt.Sprint(result))
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
