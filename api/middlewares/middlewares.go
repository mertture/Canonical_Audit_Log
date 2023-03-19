package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mertture/audit-log/api/auth"
)

func SetMiddlewareJSON(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		next(c)
	}
}

func SetMiddlewareAuthentication(next gin.HandlerFunc) gin.HandlerFunc {
    return func(c *gin.Context) {
		err := auth.TokenValid(c) // call TokenValid with the context parameter
        if err != nil {
            return
        }
        next(c)
    }
}