package main

import (
	"encoding/json"
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

	InitPgPool() // will fail the app if this doesn't work

	recipes, err := GetAllRecipes()
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(recipes)
	// output recipes as json
	json, err := json.Marshal(recipes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(json))


	router := SetupRoutes()
	fmt.Println("START API on port", AppConfig.HTTP_PORT)
	log.Fatal(http.ListenAndServe(AppConfig.HTTP_PORT, router))


}