package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kurniawanew/pfm-kambing/api/auth"
)

func SetMiddlewareAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
