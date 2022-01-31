package main

import (
	"fmt"
	"go-api-portfolio/router"
	"log"
	"net/http"
)

func main() {
	//err := godotenv.Load()

	r := router.Router()
	// fs := http.FileServer(http.Dir("build"))
	// http.Handle("/", fs)

	//s3 := os.Getenv("POSTGRES_URL")
	//fmt.Println(s3)

	fmt.Println("Starting server on the port 8080...")

	log.Fatal(http.ListenAndServe(":8080", r))
}
