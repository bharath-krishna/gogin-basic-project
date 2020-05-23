package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/coreos/go-oidc"
	"github.com/urfave/cli"
)

type Config struct {
	ServerAddr    string         `name:"addr" usage:"IP address and port to listen on" env:"ADDRESS"`
	DgraphHost    string         `name:"dgraphhost" usage:"IP address and port of dgraph" env:"DGRAPH_HOST"`
	AuthHost      string         `name:"authhost" usage:"IP address and port of auth host" env:"AUTH_HOST"`
	KeycloakRealm string         `name:"keycloak-realm" usage:"Keycloak realm used to authenticate to the oauth service" env:"KEYCLOAK_REALM"`
	ClientID      string         `name:"client_id" usage:"Client ID" env:"CLIENT_ID"`
	ClientSecret  string         `name:"client_secret" usage:"Client Secret" env:"CLIENT_SECRET"`
	Provider      *oidc.Provider `json:"provider" usage:"OIDC Provider"`
	Verbose       bool           `name:"verbose" usage:"switch on debug / verbose logging"`
}

func NewDefaultConfig() *Config {
	return &Config{
		ServerAddr:    os.Getenv("FAMILY_ADDRESS"),
		KeycloakRealm: os.Getenv("FAMILY_KEYCLOAK_REALM"),
		ClientID:      os.Getenv("FAMILY_CLIENT_ID"),
		ClientSecret:  os.Getenv("FAMILY_CLIENT_SECRET"),
		AuthHost:      os.Getenv("FAMILY_AUTH_HOST"),
		DgraphHost:    os.Getenv("FAMILY_DGRAPH_HOST"),
	}
}

func getCommandLineOptions() []cli.Flag {
	defaults := NewDefaultConfig()
	var flags []cli.Flag
	count := reflect.TypeOf(Config{}).NumField()
	for i := 0; i < count; i++ {
		field := reflect.TypeOf(Config{}).Field(i)
		usage, found := field.Tag.Lookup("usage")
		if !found {
			continue
		}
		envName := field.Tag.Get("env")
		if envName != "" {
			envName = envPrefix + envName
		}
		optName := field.Tag.Get("name")

		switch t := field.Type; t.Kind() {
		case reflect.Bool:
			dv := reflect.ValueOf(defaults).Elem().FieldByName(field.Name).Bool()
			msg := fmt.Sprintf("%s (default: %t)", usage, dv)
			flags = append(flags, cli.BoolTFlag{
				Name:   optName,
				Usage:  msg,
				EnvVar: envName,
			})
		case reflect.String:
			defaultValue := reflect.ValueOf(defaults).Elem().FieldByName(field.Name).String()
			flags = append(flags, cli.StringFlag{
				Name:   optName,
				Usage:  usage,
				EnvVar: envName,
				Value:  defaultValue,
			})
		}
	}

	return flags
}

func parseCLIOptions(ctx *cli.Context, config *Config) (err error) {
	// iterate the Config and grab command line options via reflection
	count := reflect.TypeOf(config).Elem().NumField()
	for i := 0; i < count; i++ {
		field := reflect.TypeOf(config).Elem().Field(i)
		name := field.Tag.Get("name")

		if ctx.IsSet(name) {
			switch field.Type.Kind() {
			case reflect.Bool:
				reflect.ValueOf(config).Elem().FieldByName(field.Name).SetBool(ctx.Bool(name))
			case reflect.String:
				reflect.ValueOf(config).Elem().FieldByName(field.Name).SetString(ctx.String(name))
			}
		}
	}
	return nil
}
