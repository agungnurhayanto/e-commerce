package main

import (
	"e-commerce/internal/config"
	"fmt"
	"log"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Application:", cfg.AppName)
	fmt.Println("Port: ", cfg.AppPort)
	fmt.Println("Database:", cfg.DBName)
}
