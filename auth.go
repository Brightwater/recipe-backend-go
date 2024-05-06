package main

import (
	"fmt"
	"io"
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
func VerifyTokenAndScope(token string) (string, error) {

	token = strings.TrimPrefix(token, "Bearer ")

	username, found := tokenCache.Get(token)
	if found {
		log.Println("Token validated using cache")
		// get username as str
		usernameStr, _ := username.(string)
		usernameStr = strings.ReplaceAll(usernameStr, "\"", "")
		fmt.Println(usernameStr)
		return usernameStr, nil
	}

	url := AppConfig.AUTH_BASE_PATH + "/verifyTokenAndScope?token=" + token + "&scope=test"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("unauthorized")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// get string from body NOT as byte array
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyStr := strings.ReplaceAll(string(bodyBytes), "\"", "")
	fmt.Println(bodyStr)

	tokenCache.Set(token, bodyStr, cache.DefaultExpiration)

	log.Println("Token validated using api")

	return bodyStr, nil
}
