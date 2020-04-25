package main

import (
	_ "family-tree/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) Routes(r *gin.Engine) {
	r.GET("/welcome", s.Welcome)
	r.GET("/welcome2", s.Welcome2)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
