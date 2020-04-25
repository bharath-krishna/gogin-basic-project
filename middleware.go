package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (s *Server) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			s.logger.Info("client request",
				zap.Duration("latency", time.Since(start)),
				zap.Int("status", c.Writer.Status()),
				zap.String("requester", c.Request.RemoteAddr),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.RequestURI))
		}()
		c.Next()
	}
}

func (s *Server) AllowOriginRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}
