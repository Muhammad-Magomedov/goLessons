package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type HttpResponse struct {
	Json Post `json:"json"`
}

func main() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}

	responses, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	postsBytes, err := io.ReadAll(responses.Body)
	if err != nil {
		fmt.Printf("Ошибка %v", err)
	}

	allPosts := []Post{}
	err = json.Unmarshal(postsBytes, &allPosts)
	if err != nil {
		fmt.Printf("Не смогли анмаршалить структуру постов %v", err)
	}

	fmt.Printf("%+v", allPosts)

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка: %v", err)
	}

	post := Post{}
	err = json.Unmarshal(responseBytes, &post)
	if err != nil {
		fmt.Printf("Не смогли анмаршалить структуру %v", err)
	}

	fmt.Printf("%+v", post)
	HomeWorkPost()
}

func HomeWorkPost() {
	const url = "https://httpbin.org/post"

	var mockPost = Post{
		UserID: 1,
		ID:     1,
		Title:  "test",
		Body:   "test123",
	}

	requestBody := strings.NewReader(`
	{
	"UserID": 1,
	"ID": 1,
	"Title": "test",
	"Body": "test123"
	}
	`)

	resp, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		fmt.Printf("Error cant request %v", err)
	}

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка: %v", err)
	}

	httpResponse := HttpResponse{}
	err = json.Unmarshal(responseBytes, &httpResponse)
	if err != nil {
		fmt.Printf("Не смогли анмаршалить структуру %v", err)
	}

	if (mockPost.ID == httpResponse.Json.ID) && (mockPost.Body == httpResponse.Json.Body) && (mockPost.Title == httpResponse.Json.Title) && (mockPost.UserID == httpResponse.Json.UserID) {
		fmt.Print("Success")
	}

	fmt.Printf("%+v", httpResponse)
}
