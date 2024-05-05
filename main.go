package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	
	err := LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	router := SetupRoutes()
	fmt.Println("START API on port", AppConfig.HTTP_PORT)
	log.Fatal(http.ListenAndServe(AppConfig.HTTP_PORT, router))


}