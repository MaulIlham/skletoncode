package main

import (
	"ELKExample/conf"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "local" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading ..env file not FOUND")
		}
	}
	conf.Init()
}
