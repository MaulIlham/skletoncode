package main

import (
	"ELKExample/conf"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading ..env file not FOUND")
	}
	conf.Init()
}
