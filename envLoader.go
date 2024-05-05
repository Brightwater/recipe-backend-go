package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HOST           string
	PORT           string
	DBUSER         string
	PASSWORD       string
	DBNAME         string
	AUTH_BASE_PATH string
	HTTP_PORT      string
}

// store instance of config here
var AppConfig *Config

func LoadConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file, proceeding with environment variables if present.")
		return err
	}

	log.Println("Loading .env")

	AppConfig = &Config{
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("DBUSER"),
		os.Getenv("PASSWORD"),
		os.Getenv("DBNAME"),
		os.Getenv("AUTH_BASE_PATH"),
		":" + os.Getenv("HTTP_PORT"),
	}

	return nil
}
