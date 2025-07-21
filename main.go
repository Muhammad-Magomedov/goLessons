package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
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

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка: %v", err)
	}

	post := Post{}
	err = json.Unmarshal(bytes, &post)
	if err != nil {
		fmt.Printf("Не смогли анмаршалить структуру %v", err)
	}

	fmt.Printf("%+v", post)
}
