package main

import (
	"fmt"
	"github.com/samuelg/rentals/config"
	"os"
)

func getEnv() string {
	if value, ok := os.LookupEnv("ENV"); ok {
		return value
	}
	// default to development
	return "development"
}

func main() {
	env := getEnv()
	fmt.Printf("Loanding config for %s environment\n", env)
	config.Init(env)

	fmt.Println("Rental API service")
}
