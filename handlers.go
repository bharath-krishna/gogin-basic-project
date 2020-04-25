package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Welcome godoc
// @Summary welcome page
// @Description This is the welcome page
// @Tags welcome
// @Accept  json
// @Produce  json
// @Produce  json
// @Success 200 {object} Welcome
// @Router /welcome [get]
func (s *Server) Welcome(c *gin.Context) {
	ret := Welcome{Id: "1"}
	c.JSON(http.StatusOK, ret)
}

type Welcome struct {
	Id string `json:"id" example:"1" format:"int64"`
}

// Welcome2 godoc
// @Summary welcome2 page
// @Description This is the welcome2 page
// @Tags welcome2
// @Accept  json
// @Produce  json
// @Produce  text/plain
// @Success 200 {string} string
// @Success 400 {object} Welcome2
// @Router /welcome2 [post]
func (s *Server) Welcome2(c *gin.Context) {
	ret := "srtfsdfsd"
	c.JSON(http.StatusOK, ret)
}

type Welcome2 struct {
	Name string `json:"name" example:"somename" format:"string"`
}
