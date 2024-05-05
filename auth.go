package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

)

var tokenCache = cache.New(15*time.Minute, 15*time.Minute)

// check the auth header and see if it contains a bearer token
// func VerifyAuth(r *http.Request) error {
// 	token := r.Header.Get("Authorization")
// 	if token == "" {
// 		return fmt.Errorf("no token")
// 	}

// 	log.Printf("Auth header received: %s", token)

// 	if strings.Contains(token, "Bearer") {
// 		return VerifyTokenAndScope(r)
// 	}

// 	return fmt.Errorf("Auth failed")
// }

// call the oauth service and check the token
func VerifyTokenAndScope(token string) error {

	token = strings.TrimPrefix(token, "Bearer ")

	_, found := tokenCache.Get(token)
	if found {
		log.Println("Token validated using cache")
		return nil
	}

	url := AppConfig.AUTH_BASE_PATH + "/verifyTokenAndScope?token=" + token + "&scope=test"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("unauthorized")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	tokenCache.Set(token, token, cache.DefaultExpiration)

	return nil
}
