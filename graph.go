package main

import (
	"context"
	"strings"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// func getNewGraphClient(logger *zap.Logger, config *Config) (*dgo.Dgraph, error) {
// 	var clients []api.DgraphClient
// 	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
// 	if err != nil {
// 		fmt.Printf("********************%+v********************\n", "error is")
// 		fmt.Printf("********************%+v********************\n", err)
// 	}
// 	defer conn.Close()

// 	dc := api.NewDgraphClient(conn)
// 	clients = append(clients, dc)
// 	gclient := dgo.NewDgraphClient(clients...)
// 	return gclient, nil
// }

type Client struct {
	dgraph *dgo.Dgraph
	logger *zap.Logger
	config *Config
}

func GetNewGraphClient(logger *zap.Logger, config *Config) (*Client, error) {
	var clients []api.DgraphClient
	for _, d := range strings.Split(config.DgraphHost, ",") {
		conn, err := grpc.Dial(d, grpc.WithInsecure())
		if err != nil {
			logger.Error(
				"Upnable to create connection",
				zap.Error(err),
			)
			return nil, err
		}
		client := api.NewDgraphClient(conn)
		clients = append(clients, client)
	}
	dgraph := dgo.NewDgraphClient(clients...)

	resp, err := dgraph.NewTxn().Query(context.Background(), `schema{}`)
	if err != nil {
		logger.Fatal(
			"Unable to retrive schema",
			zap.Error(err),
		)
		return nil, err
	}
	logger.Debug("Recived schema from dgraph", zap.String("schema", string(resp.Json)))

	return &Client{dgraph: dgraph, logger: logger, config: config}, nil
}
