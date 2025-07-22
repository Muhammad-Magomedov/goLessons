package main

import (
	"bytes"
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

type HttpResponse struct {
	Json Post `json:"json"`
}

func main() {
	const url = "https://httpbin.org/post"

	var mockPost = Post{
		UserID: 1,
		ID:     1,
		Title:  "test",
		Body:   "test123",
	}

	jsonData, err := json.Marshal(mockPost)
	if err != nil {
		fmt.Printf("Error, cant marshal data to json %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
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
