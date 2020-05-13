package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/dgo/v2/protos/api"
)

type Person struct {
	Name      string    `json:"name"`
	ID        string    `json:"id,omitempty"`
	Father    *Person   `json:"father,omitempty"`
	Mother    *Person   `json:"mother,omitempty"`
	Sons      []*Person `json:"sons,omitempty"`
	Wife      []*Person `json:"wife,omitempty"`
	Husband   []*Person `json:"husband,omitempty"`
	Daughters []*Person `json:"daughters,omitempty"`
	Gender    string    `json:"gender,omitempty"`
	Deleted   bool      `json:"deleted,omitempty"`
	Deceased  bool      `json:"deceased,omitempty"`
	DType     []string  `json:"dgraph.type,omitempty"`
}

var (
	SEARCH_CREATE_QUERY = `upsert {
		query {
			var(func: eq(name, "$father_name")) {
				father as uid
				partner @filter(eq(name, "$mother_name")) {
					mother as uid
				}
			}
		}
		mutation {
			set {
				_:$name <dgraph.type> "Person" .
				_:$name <name> "$name" .
				_:$name <father> uid(father) .
				_:$name <mother> uid(mother) .
				_:$name <gender> "$gender" .
			}
		}
	}`

	SEARCH_QUERY_BY_NAME = `{
		person(func: eq(name, "%s")) {
		  id: uid
		  name
		  wife {
			  id: uid 
			  name
			  gender
		  }
		  husband {
			  id: uid 
			  name
			  gender
		  }
		  father {
			  name
			  gender
		  }
		  mother {
			  name
			  gender
		  }
		  gender
		}
	  }`
	SEARCH_QUERY_BY_ID = `{
		person(func: uid(%s)) {
		  id: uid
		  name
		  wife {
			  id: uid
			  name
			  gender
		  }
		  husband {
			  id: uid
			  name
			  gender
		  }
		  father {
			  id: uid
			  name
			  gender
		  }
		  mother {
			  id: uid
			  name
			  gender
		  }
		  gender
		}
	  }`
	QUERY_CHILDREN = `{
		person(func: uid(%s)) {
			id: uid
			name
			sons: ~%s @filter(eq(gender, "male")) {
				id: uid
				name
			}
			daughters: ~%s @filter(eq(gender, "female")) {
				id: uid
				name
			}
		}
	}`
	QUERY_HUSBAND_OR_WIFE = `{
		person(func: uid(%s)) {
			id: uid
			name
			gender
			wife {
				id: uid
				name
				gender
			}
			husband {
				id: uid
				name
				gender
			}
		}
	}`
	QUERY_FATHER = `{
		person(func: uid(%s)) {
			id: uid
			name
			gender
			father {
				id: uid
				name
				gender
			}
		}
	}`
	QUERY_MOTHER = `{
		person(func: uid(%s)) {
			id: uid
			name
			gender
			mother {
				id: uid
				name
				gender
			}
		}
	}`
	QUERY_ALL = `{
		person(func: has(name)) @recurse(depth: 2, loop: true) {
			id: uid
			expand(_all_)
		}
	}`
	QUERY_ALL_NETWORK_FORMAT = `{
		nodes(func: has(name)) {
		  id: uid
		  expand(_all_)
		}
		fathers(func: has(name)) @cascade @normalize {
		  source: uid
		  sname: name
		  father {
			target: uid
			tname: name
		  }
		}
		mothers(func: has(name)) @cascade @normalize {
		  source: uid
		  sname: name
		  mother {
			target: uid
			tname: name
		  }
		}
		wife(func: has(name)) @cascade @normalize {
		  source: uid
		  sname: name
		  husband {
			target: uid
			tname: name
		  }
		}
		husband(func: has(name)) @cascade @normalize {
			source: uid
			sname: name
			wife {
			  target: uid
			  tname: name
			}
		  }
	  }`
	QUERY_PERSON_NETWORK_FORMAT = `{
		nodes(func: has(name)) @filter(uid_in(~father, %s) OR uid_in(father, %s) OR uid_in(mother, %s) OR uid_in(~wife, %s) OR uid_in(~mother, %s) OR uid_in(~husband, %s) OR uid(%s)) {
		  id: uid
          name
          gender
		}
		fathers(func: uid(%s)) @cascade @normalize {
		  source: uid
		  sname: name
		  father {
			target: uid
			tname: name
		  }
		}
		mothers(func: uid(%s)) @cascade @normalize {
		  source: uid
		  sname: name
		  mother {
		    target: uid
			tname: name
		  }
		}
		wife(func: uid(%s)) @cascade @normalize {
		  source: uid
		  sname: name
		  husband {
		    target: uid
			tname: name
		  }
		}
		husband(func: uid(%s)) @cascade @normalize {
		  source: uid
		  sname: name
		  wife {
			target: uid
			tname: name
		  }
		}
		fsons(func: uid(%s)) @cascade @normalize {
          source: uid
		  sname: name
          ~father @filter(eq(gender, "male")) {
            target: uid
			tname: name
          }
        }
		msons(func: uid(%s)) @cascade @normalize {
            source: uid
			sname: name
            ~mother @filter(eq(gender, "male")) {
              target: uid
			  tname: name
            }
          }
          fdaughters(func: uid(%s)) @cascade @normalize {
            source: uid
			sname: name
            ~father @filter(eq(gender, "female")) {
              target: uid
			  tname: name
            }
          }
          mdaughters(func: uid(%s)) @cascade @normalize {
            source: uid
			sname: name
            ~mother @filter(eq(gender, "female")) {
			  target: uid
			  tname: name
            }
          }
        }`
)

func (c *Client) CreatePerson(p *Person) error {
	p.DType = []string{"Person"}
	ctx := context.Background()
	txn := c.dgraph.NewTxn()
	defer txn.Discard(ctx)

	pb, err := json.Marshal(p)
	if err != nil {
		c.logger.Fatal(err.Error())
	}

	mu := &api.Mutation{
		SetJson:   pb,
		CommitNow: true,
	}
	_, err = txn.Mutate(ctx, mu)
	if err != nil {
		c.logger.Fatal(err.Error())
	}
	return nil
}

func (c *Client) DeletePerson(id string) error {
	p := &Person{}
	p.ID = id
	p.DType = []string{"Person"}
	ctx := context.Background()
	txn := c.dgraph.NewTxn()
	defer txn.Discard(ctx)

	// Upsert query
	// upsert {
	// 	query {
	// 	  p as node(func: id(0x7566)){
	// 		h as ~husband
	// 		w as ~wife
	// 		f as ~father
	// 		m as ~mother
	// 	  }
	// 	}

	// 	mutation {
	// 	  delete {
	// 		uid(h) <husband> uid(p) .
	// 		uid(w) <wife> uid(p) .
	// 		uid(f) <father> uid(p) .
	// 		uid(m) <mother> uid(p) .
	// 		uid(p) * * .
	// 	  }
	// 	}
	//   }

	query := fmt.Sprintf(`query {
		person as node(func: uid(%s)){
			h as ~husband
			w as ~wife
			f as ~father
			m as ~mother	  
		}
	}`, id)

	mu := &api.Mutation{
		DelNquads: []byte(`
		uid(h) <husband> uid(person) .
		uid(w) <wife> uid(person) .
		uid(f) <father> uid(person) .
		uid(m) <mother> uid(person) .
		uid(person) <name> * .
		uid(person) <gender> * .
		uid(person) <father> * .
		uid(person) <mother> * .
		uid(person) * * .
		  `),
	}
	req := &api.Request{
		Query:     query,
		Mutations: []*api.Mutation{mu},
		CommitNow: true,
	}

	// Update email only if matching uid found.
	if _, err := c.dgraph.NewTxn().Do(ctx, req); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (c *Client) SearchPerson(query string) ([]Person, error) {
	ctx := context.Background()
	txn := c.dgraph.NewTxn()
	defer txn.Discard(ctx)

	req := &api.Request{
		Query: query,
	}
	res, err := txn.Do(ctx, req)
	if err != nil {
		c.logger.Fatal(err.Error())
	}

	reqData := map[string][]Person{}
	if err := json.Unmarshal(res.Json, &reqData); err != nil {
		c.logger.Fatal(err.Error())
	}

	return reqData["person"], nil
}

func (c *Client) GetPropleNetwork(query string) (map[string][]map[string]string, error) {
	ctx := context.Background()
	txn := c.dgraph.NewTxn()
	defer txn.Discard(ctx)

	req := &api.Request{
		Query: query,
	}
	res, err := txn.Do(ctx, req)
	if err != nil {
		c.logger.Fatal(err.Error())
	}

	reqData := map[string][]map[string]string{}
	if err := json.Unmarshal(res.Json, &reqData); err != nil {
		c.logger.Fatal(err.Error())
	}

	return reqData, nil
}

func (c *Client) UpdatePerson(p *Person) error {
	p.DType = []string{"Person"}
	ctx := context.Background()
	txn := c.dgraph.NewTxn()
	defer txn.Discard(ctx)

	pb, err := json.Marshal(p)
	if err != nil {
		c.logger.Fatal(err.Error())
	}

	mu := &api.Mutation{
		SetJson:   pb,
		CommitNow: true,
	}
	_, err = txn.Mutate(ctx, mu)
	if err != nil {
		c.logger.Fatal(err.Error())
	}
	return nil
}

// func (c *Client) UpdatePartners(p *Person) error {
// 	p.DType = []string{"Person"}
// 	query := fmt.Sprintf(SEARCH_QUERY_BY_UID, p.UID)
// 	person, err := c.SearchPerson(query)
// 	if err != nil {
// 		c.logger.Fatal(err.Error())
// 		return (err)
// 	}
// 	p.Gender = person[0].Gender
// 	p.Name = person[0].Name
// 	if person[0].Gender == "male" {
// 		p.Partner = person[0].Wife
// 	}
// 	ctx := context.Background()
// 	txn := c.dgraph.NewTxn()
// 	defer txn.Discard(ctx)

// 	pb, err := json.Marshal(p)
// 	if err != nil {
// 		c.logger.Fatal(err.Error())
// 	}

// 	mu := &api.Mutation{
// 		SetJson:   pb,
// 		CommitNow: true,
// 	}
// 	_, err = txn.Mutate(ctx, mu)
// 	if err != nil {
// 		c.logger.Fatal(err.Error())
// 	}
// 	return nil
// }
