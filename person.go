package main

import (
	"context"
	"encoding/json"

	"github.com/dgraph-io/dgo/v2/protos/api"
)

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
		  uid
		  name
		  partner {
			  uid
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
	SEARCH_QUERY_BY_UID = `{
		person(func: uid("%s")) {
		  uid
		  name
		  partner {
			  uid
			  name
			  gender
		  }
		  father {
			  uid
			  name
			  gender
		  }
		  mother {
			  uid
			  name
			  gender
		  }
		  gender
		}
	  }`
	DELETE_QUERY_BY_UID = `{
		person(func: uid("%s")) {
		  uid
		  name
		  partner {
			  uid
			  name
			  gender
		  }
		  father {
			  uid
			  name
			  gender
		  }
		  mother {
			  uid
			  name
			  gender
		  }
		  gender
		}
	  }`
	QUERY_CHILDREN = `{
		person(func: uid(%s)) {
			uid
			name
			sons: ~father @filter(eq(gender, "male")) {
				uid
				name
			}
			daughters: ~father @filter(eq(gender, "female")) {
				uid
				name
			}
		}
	}`
	QUERY_PARTNERS = `{
		person(func: uid(%s)) {
			uid
			name
			gender
			wife: partner {
				uid
				name
				gender
			}
			husband: ~partner {
				uid
				name
				gender
			}
		}
	}`
	QUERY_FATHER = `{
		person(func: uid(%s)) {
			uid
			name
			gender
			father {
				uid
				name
				gender
			}
		}
	}`
	QUERY_MOTHER = `{
		person(func: uid(%s)) {
			uid
			name
			gender
			mother {
				uid
				name
				gender
			}
		}
	}`
	QUERY_ALL = `{
		person(func: has(name)) @recurse(depth: 2, loop: true) {
			uid
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
		  father {
			target: uid
		  }
		}
		mothers(func: has(name)) @cascade @normalize {
		  source: uid
		  mother {
			target: uid
		  }
		}
		partners(func: has(name)) @cascade @normalize {
		  source: uid
		  partner {
			target: uid
		  }
		}
	}`
)

type Person struct {
	Name      string    `json:"name"`
	UID       string    `json:"uid,omitempty"`
	Partner   []*Person `json:"partner,omitempty"`
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
