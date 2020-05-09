package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPerson godoc
// @Summary GetPerson
// @Description Get details of a person like father, mother and partners via uid
// @Tags Person
// @Accept  json
// @Produce  json
// @Success 200 {object} Person
// @Router /{uid}/ [get]
// @Param uid path string true "uid"
func (s *Server) GetPerson(c *gin.Context) {
	person := &Person{UID: c.Param("uid")}
	query := fmt.Sprintf(SEARCH_QUERY_BY_UID, person.UID)

	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	c.JSON(http.StatusOK, people[0])
}

func (s *Server) DeletePerson(c *gin.Context) {
	person := &Person{UID: c.Param("uid")}
	query := fmt.Sprintf(SEARCH_QUERY_BY_UID, person.UID)

	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	c.JSON(http.StatusOK, people[0])
}

// GetChildren godoc
// @Summary GetChildren
// @Description Get children's details of a person via uid (only immediage children are fetched)
// @Tags Person
// @Accept  json
// @Produce  json
// @Success 200 {object} Person
// @Router /{uid}/children/ [get]
// @Param uid path string true "uid"
func (s *Server) GetChildren(c *gin.Context) {
	person := &Person{UID: c.Param("uid")}
	query := fmt.Sprintf(QUERY_CHILDREN, person.UID)

	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if len(people) > 0 {
		c.JSON(http.StatusOK, people[0])
	} else {
		c.JSON(http.StatusOK, people)
	}
}

func (s *Server) GetPartners(c *gin.Context) {
	person := &Person{UID: c.Param("uid")}
	query := fmt.Sprintf(QUERY_PARTNERS, person.UID)

	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if len(people) > 0 {
		c.JSON(http.StatusOK, people[0])
	} else {
		c.JSON(http.StatusOK, people)
	}
}

func (s *Server) GetFather(c *gin.Context) {
	person := &Person{UID: c.Param("uid")}
	query := fmt.Sprintf(QUERY_FATHER, person.UID)

	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if len(people) > 0 {
		c.JSON(http.StatusOK, people[0])
	} else {
		c.JSON(http.StatusOK, people)
	}
}

func (s *Server) GetMother(c *gin.Context) {
	person := &Person{UID: c.Param("uid")}
	query := fmt.Sprintf(QUERY_MOTHER, person.UID)

	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if len(people) > 0 {
		c.JSON(http.StatusOK, people[0])
	} else {
		c.JSON(http.StatusOK, people)
	}
}

// UpdatePerson godoc
// @Summary UpdatePerson
// @Description Update person's details like father, mother or partners by post data, uid of a person is required in url path
// @Tags Person
// @Accept  json
// @Produce  json
// @Success 200 {object} Person
// @Router /{uid}/ [patch]
// @Param uid path string true "uid"
// @Param person body Person true "json data"
func (s *Server) UpdatePerson(c *gin.Context) {
	patchPerson := c.MustGet("person").(*Person)
	patchPerson.UID = c.Param("uid")
	err := s.gclient.UpdatePerson(patchPerson)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, patchPerson)
}

// CreatePerson godoc
// @Summary CreatePerson
// @Description Create a new entry for a person object, you can mention any level of father, mother or partners.
// @Tags People
// @Accept  json
// @Produce  json
// @Success 200 {object} Person
// @Router / [post]
// @Param person body Person true "json data"
func (s *Server) CreatePerson(c *gin.Context) {
	person := c.MustGet("person").(*Person)
	if err := s.gclient.CreatePerson(person); err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, person)
}

// SearchPerson godoc
// @Summary SearchPerson
// @Description Search for a person by the given fields, all the matching objects will be returned as a list
// @Tags People
// @Accept  json
// @Produce  json
// @Success 200 {object} Person
// @Router /search [post]
// @Param person body Person true "json data"
func (s *Server) SearchPerson(c *gin.Context) {
	searchPerson := c.MustGet("person").(*Person)
	query := fmt.Sprintf(SEARCH_QUERY_BY_NAME, searchPerson.Name)
	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"people": people})
}

// GetAllPeople godoc
// @Summary GetALlPeople
// @Description Get all people data, this includes self and one level like parent's and partner's data
// @Tags People
// @Accept  json
// @Produce  json
// @Success 200 {object} Person
// @Router / [get]
func (s *Server) GetAllPeople(c *gin.Context) {
	people, err := s.gclient.SearchPerson(QUERY_ALL)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"people": people})
}

func (s *Server) GetAllPeopleNetwork(c *gin.Context) {
	data, err := s.gclient.GetPropleNetwork(QUERY_ALL_NETWORK_FORMAT)
	links := append(data["mothers"], data["fathers"]...)
	links = append(links, data["partners"]...)
	peopleNetwork := map[string][]map[string]string{"nodes": data["nodes"], "links": links}
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": peopleNetwork})
}
