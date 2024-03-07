package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Company         string `json:"company"`
	Location        string `json:"location"`
	Email           string `json:"email"`
	Bio             string `json:"bio"`
	TwitterUsername string `json:"twitter_username"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	UserUrl         string `json:"url"`
}

var users []User

func main() {
	fmt.Println("hello world")

	var wg sync.WaitGroup

	if len(users) == 0 {
		getAllUsers()

		fmt.Println(len(users))

		waitCounter := len(users)
		wg.Add(waitCounter)
		for u := 0; u < len(users); u++ {
			go getUserDetails(u, &wg)
		}
		wg.Wait()
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": users[:10],
		})
	})

	r.Run()

}

func getAllUsers() {
	log.Println("*****************inside get all users")
	client := &http.Client{}

	// resp, err := http.Get("https://api.github.com/users")
	req, err := http.NewRequest("GET", "https://api.github.com/users", nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", os.Getenv("GITHUB_TOKEN"))

	resp, err := client.Do(req)
	fmt.Println("***********************1*****************************")
	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println("***********************2*****************************")

	// Unmarshal JSON data into a slice of items
	// var items []User
	err = json.Unmarshal(body, &users)
	fmt.Println("***********************3*****************************")

	if err != nil {
		var jsonData interface{}
		log.Print(json.Unmarshal(body, &jsonData))
		log.Fatalln(err)
		return
	}
	fmt.Println("***********************4*****************************")

	fmt.Println(users)

	return
}

func getUserDetails(u int, wg *sync.WaitGroup) {
	// fmt.Println(users[u].UserUrl, "userurl************")
	client := &http.Client{}

	req, err := http.NewRequest("GET", users[u].UserUrl, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", os.Getenv("GITHUB_TOKEN"))

	resp, err := client.Do(req)

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}

	var userObj User
	json.Unmarshal(body, &userObj)

	users[u].Name = userObj.Name
	users[u].Company = userObj.Company
	users[u].Location = userObj.Location
	users[u].Email = userObj.Email
	users[u].Bio = userObj.Bio
	users[u].TwitterUsername = userObj.TwitterUsername
	users[u].CreatedAt = userObj.CreatedAt
	users[u].UpdatedAt = userObj.UpdatedAt

	log.Print(users[u])

	wg.Done()
}
