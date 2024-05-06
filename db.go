package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Recipe struct {
	ID     int         `json:"id"`
	Author string      `json:"author"`
	Data   interface{} `json:"data"`
	TS     *time.Time  `json:"ts"`
}

var Pool *pgxpool.Pool

// cache array of recipes
var allRecipes *[]Recipe
var recipesMutex sync.Mutex

func InitPgPool() {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		AppConfig.HOST, AppConfig.PORT, AppConfig.DBUSER, AppConfig.PASSWORD, AppConfig.DBNAME)

	poolConfig, err := pgxpool.ParseConfig(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	poolConfig.MinConns = 0
	poolConfig.MaxConns = 3
	poolConfig.MaxConnIdleTime = 15 * time.Minute
	poolConfig.MaxConnLifetime = 30 * time.Minute

	Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Database connected")
}

func GetAllRecipes() (*[]Recipe, error) {

	// check len of allRecipes
	if allRecipes != nil && len(*allRecipes) > 0 {
		return allRecipes, nil
	}

	conn, err := Pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	queryStry := "SELECT id, author, data, ts FROM recipes"

	result, err := conn.Query(context.Background(), queryStry)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	recipes := make([]Recipe, 0)
	// Loop through rows, using Scan to assign column data to struct fields
	for result.Next() {
		recipe := new(Recipe)
		err = result.Scan(&recipe.ID, &recipe.Author, &recipe.Data, &recipe.TS)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, *recipe)
	}

	// cache recipes
	recipesMutex.Lock()
	defer recipesMutex.Unlock()
	allRecipes = &recipes

	return &recipes, nil
}

func UpdateRecipe(recipe *Recipe) error {

	conn, err := Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	updateStr := "UPDATE recipes SET data = $1 WHERE id = $2 AND author = $3"

	_, err = conn.Exec(context.Background(), updateStr, recipe.Data, recipe.ID, recipe.Author)
	if err != nil {
        fmt.Println("Db error adding recipe in UpdateRecipe", err)
		return err
	}

	// evict recipes cache
	recipesMutex.Lock()
	defer recipesMutex.Unlock()
	allRecipes = nil

	return nil
}

func AddRecipe(recipe *Recipe) error {

	conn, err := Pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	insertStr := "INSERT INTO recipes (author, data, ts) VALUES ($1, $2, now())"

	_, err = conn.Exec(context.Background(), insertStr, recipe.Author, recipe.Data)
	if err != nil {
		return err
	}

	// evict recipes cache
	recipesMutex.Lock()
	defer recipesMutex.Unlock()
	allRecipes = nil

	return nil
}
