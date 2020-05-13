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
	person.GET("/:id", s.GetPerson)
	person.GET("/:id/network", s.GetPersonNetwork)
	person.DELETE("/:id", s.DeletePerson)
	person.GET("/:id/children", s.GetChildren)
	person.GET("/:id/husband", s.GetHusband)
	person.GET("/:id/wife", s.GetWife)
	person.GET("/:id/father", s.GetFather)
	person.GET("/:id/mother", s.GetMother)

	person.Use(s.FetchPerson())
	{
		// .../people/* endpooint
		person.POST("/", s.CreatePerson)
		person.POST("/search", s.SearchPerson)
		person.PATCH("/:id", s.UpdatePerson)
	}
}
