package errors

import "github.com/gin-gonic/gin"

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(
		code,
		gin.H{
			"status": false,
			"error":  message,
		},
	)
}
