package main

import (
	"github.com/joho/godotenv"
	"kpdigisign/routes"
	"os"
	"path/filepath"
)

func init()  {
	fileExecutable, _ := os.Executable()
	basePath, _ := filepath.Split(fileExecutable)
	if os.Getenv("APP_ENV") != "production" {
		basePath = ""
	}
	_ = godotenv.Load(basePath + ".env")
}

func main() {
	e :=routes.New()

	e.Logger.Fatal(e.Start(":"+os.Getenv("APP_PORT")))
}


