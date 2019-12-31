package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"kpdigisign/routes"
)

func init()  {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {
	e :=routes.New()

	e.Logger.Fatal(e.Start(":1323"))
}


