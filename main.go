package main

import (
	_ "family-tree/docs"
	"os"
)

// @title Family tree apis
// @version 3.0.0
// @description This is api spec for family tree project.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8088
// @BasePath /people
func main() {
	app := newApp()
	app.Run(os.Args)
}
