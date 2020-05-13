package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		c.Header("Access-Control-Allow-Origin", "*")
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

func (s *Server) FetchPerson() gin.HandlerFunc {
	return func(c *gin.Context) {
		person, err := getBody(c)
		if err != nil {
			fmt.Printf("********************%+v********************\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			c.Set("person", person)
			c.Next()
		}
	}
}

func getBody(c *gin.Context) (*Person, error) {
	reqData := &Person{}
	if profBytes, err := ioutil.ReadAll(c.Request.Body); err == nil {
		if err := json.Unmarshal(profBytes, reqData); err == nil {
			return reqData, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
