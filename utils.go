package main

import (
	"context"
	"fmt"
)

func (c *Client) GetUID(p *Person) (map[string]string, error) {

	q := `query all($a: string) {
		all(func: eq(name, $a)) {
		  name
		}
	}`

	ctx := context.Background()
	txn := c.dgraph.NewTxn()
	defer txn.Discard(ctx)

	res, err := txn.QueryWithVars(ctx, q, map[string]string{"$a": p.Name})
	fmt.Printf("********************%+v********************\n", res)
	if err != nil {
		fmt.Printf("********************%+v********************\n", err)
		return nil, err
	}
	fmt.Printf("%s\n", res.Json)
	return res.Uids, nil
}
