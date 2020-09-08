package main

import (
	_ "family-tree/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) Routes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/network", s.GetAllPeopleNetwork)
	person := r.Group("/people")
	person.GET("/", s.GetAllPeople)
	person.GET("/:uid", s.GetPerson)
	person.GET("/:uid/network", s.GetPersonNetwork)
	person.DELETE("/:uid", s.DeletePerson)
	person.GET("/:uid/children", s.GetChildren)
	person.GET("/:uid/husband", s.GetHusband)
	person.GET("/:uid/wife", s.GetWife)
	person.GET("/:uid/father", s.GetFather)
	person.GET("/:uid/mother", s.GetMother)

	person.Use(s.FetchPerson())
	{
		// .../people/* endpooint
		person.POST("/", s.CreatePerson)
		person.POST("/search", s.SearchPerson)
		person.PATCH("/:uid", s.UpdatePerson)
	}
}
