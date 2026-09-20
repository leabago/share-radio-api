package main

import (
	"fmt"
	"log"

	"github.com/leabago/share-radio/adder/config"
	"github.com/leabago/share-radio/adder/internal/app"
)

func main() {
	fmt.Println("start 1")
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
