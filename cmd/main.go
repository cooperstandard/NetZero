package main

import (
	"fmt"

	"github.com/cooperstandard/NetZero/internal/database"
)

type apiConfig struct {
	db          *database.Queries
	tokenSecret string
}

func main() {
	fmt.Println("Welcome")
}
