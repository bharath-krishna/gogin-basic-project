package main

import (
	"os"
)

func main() {
	app := newApp()
	app.Run(os.Args)
}
