package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

}

func main() {
	local_db := os.Getenv("LOCAL_DB")
	fmt.Println("Application Run ...", local_db)
}
