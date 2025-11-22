package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReadOnlyMiddleware(isReadOnly bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		if isReadOnly {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Read-only replica",
				"message": "Write operations not allowed on read-only instance",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func InstanceInfoMiddleware(port string, instanceName string, isReadOnly bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Backend-Port", port)
		c.Header("X-Backend-Mode", map[bool]string{true: "RO", false: "RW"}[isReadOnly])
		c.Header("X-Backend-Instance", instanceName)
		c.Next()
	}
}
