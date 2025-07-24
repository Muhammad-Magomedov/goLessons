package main

import (
	"fmt"
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
	res := isPalindrome(10)
	fmt.Println(res)
}

func isPalindrome(x int) bool {
	result := reverseConcat(strings.Join(strings.Split(string(x), ""), ""))
	if result == strings.Join(strings.Split(string(x), ""), "") {
		return true
	}
	return false
}

func reverseConcat(str string) string {
	var reversed string

	for i := len(str) - 1; i >= 0; i-- {
		reversed += string(str[i])
	}

	return reversed
}
