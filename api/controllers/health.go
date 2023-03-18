package controllers

import (
	"context"
	"net/http"
	"github.com/gin-gonic/gin"
)

func (server *Server) HealthCheck(c *gin.Context) {
    if err := server.DB.Client().Ping(context.Background(), nil); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database connection error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
