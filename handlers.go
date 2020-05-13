package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPerson godoc
// @Summary GetPerson
// @Description Get details of a person like father, mother and partners via id
// @Tags Person
// @Accept  json
// @Produce  json
// @Success 200 {object} Person
// @Router /{id}/ [get]
// @Param id path string true "id"
func (s *Server) GetPerson(c *gin.Context) {
	person := &Person{ID: c.Param("id")}
	query := fmt.Sprintf(SEARCH_QUERY_BY_ID, person.ID)

	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	c.JSON(http.StatusOK, people[0])
}

// GetChildren godoc
// @Summary GetChildren
// @Description Get children's details of a person via id (only immediage children are fetched)
// @Tags Person
// @Accept  json
// @Produce  json
// @Success 200 {object} Person
// @Router /{id}/children/ [get]
// @Param id path string true "id"
func (s *Server) GetChildren(c *gin.Context) {
	person := &Person{ID: c.Param("id")}

	personQuery := fmt.Sprintf(SEARCH_QUERY_BY_ID, person.ID)

	querPerson, err := s.gclient.SearchPerson(personQuery)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}

	parent := ""
	if querPerson[0].Gender == "female" {
		parent = "mother"
	} else {
		parent = "father"
	}

	query := fmt.Sprintf(QUERY_CHILDREN, person.ID, parent, parent)

	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if len(people) == 1 {
		c.JSON(http.StatusOK, gin.H{"sons": people[0].Sons, "daughters": people[0].Daughters})
	} else {
		children := []map[string][]*Person{}
		for _, child := range people {
			children = append(children, map[string][]*Person{"sons": child.Sons, "daughters": child.Daughters})
		}
		c.JSON(http.StatusOK, children)
	}
}

func (s *Server) GetHusband(c *gin.Context) {
	person := &Person{ID: c.Param("id")}
	query := fmt.Sprintf(QUERY_HUSBAND_OR_WIFE, person.ID)

	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if people[0].Gender == "male" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("A Male person '%s' can not have husband", people[0].Name)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"husband": people[0]})
}

func (s *Server) GetWife(c *gin.Context) {
	person := &Person{ID: c.Param("id")}
	query := fmt.Sprintf(QUERY_HUSBAND_OR_WIFE, person.ID)

	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	if people[0].Gender == "female" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("A Female person '%s' can not have wife", people[0].Name)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"wife": people[0]})
}

func (s *Server) GetFather(c *gin.Context) {
	person := &Person{ID: c.Param("id")}
	query := fmt.Sprintf(QUERY_FATHER, person.ID)

	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	c.JSON(http.StatusOK, gin.H{"father": people[0].Father})
}

func (s *Server) GetMother(c *gin.Context) {
	person := &Person{ID: c.Param("id")}
	query := fmt.Sprintf(QUERY_MOTHER, person.ID)

	people, err := s.gclient.SearchPerson(query)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	c.JSON(http.StatusOK, gin.H{"mother": people[0].Mother})
}

// UpdatePerson godoc
// @Summary UpdatePerson
// @Description Update person's details like father, mother or partners by post data, id of a person is required in url path
// @Tags Person
// @Accept  json
// @Produce  json
// @Success 200 {object} Person
// @Router /{id}/ [patch]
// @Param id path string true "id"
// @Param person body Person true "json data"
func (s *Server) UpdatePerson(c *gin.Context) {
	patchPerson := c.MustGet("person").(*Person)
	patchPerson.ID = c.Param("id")
	err := s.gclient.UpdatePerson(patchPerson)
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, patchPerson)
}

// func (s *Server) UpdatePartners(c *gin.Context) {
// 	patchPerson := c.MustGet("person").(*Person)
// 	patchPerson.id = c.Param("id")
// 	err := s.gclient.UpdatePartners(patchPerson)
// 	if err != nil {
// 		s.logger.Fatal(err.Error())
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}
// 	c.JSON(http.StatusOK, patchPerson)
// }

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

func (s *Server) DeletePerson(c *gin.Context) {
	id := c.Param("id")
	if err := s.gclient.DeletePerson(id); err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Success": fmt.Sprintf("Person with id %s has been deleted", id)})
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
	links = append(links, data["wife"]...)
	links = append(links, data["husband"]...)
	peopleNetwork := map[string][]map[string]string{"nodes": data["nodes"], "links": links}
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": peopleNetwork})
}

func (s *Server) GetPersonNetwork(c *gin.Context) {
	id := c.Param("id")
	query := fmt.Sprintf(QUERY_PERSON_NETWORK_FORMAT, id, id, id, id, id, id, id, id, id, id, id, id, id, id, id)
	data, err := s.gclient.GetPropleNetwork(query)

	// fathers := []map[string]string{}
	// for _, father := range data["fathers"] {
	// 	father["relation"] = "father"
	// 	fathers = append(fathers, father)
	// }
	// fsons := []map[string]string{}
	// for _, fson := range data["fsons"] {
	// 	fson["relation"] = "father"
	// 	fathers = append(fsons, fson)
	// }

	links := append(data["mothers"], data["fathers"]...)
	links = append(links, data["wife"]...)
	links = append(links, data["husband"]...)
	links = append(links, data["fsons"]...)
	links = append(links, data["fdaughters"]...)
	links = append(links, data["msons"]...)
	links = append(links, data["mdaughters"]...)
	peopleNetwork := map[string][]map[string]string{"nodes": data["nodes"], "links": links}
	if err != nil {
		s.logger.Fatal(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": peopleNetwork})
}
