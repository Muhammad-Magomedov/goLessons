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
	// debil := map[int]int{
	// 	1: 1,
	// }
	// res := isPalindrome(10)
	// fmt.Println(res)
	// fmt.Println(debil)

	// nums := []int{2, 7, 11, 15}
	// target := 9
	// fmt.Println(twoSum(nums, target))

	// nums = []int{3, 2, 4}
	// target = 6
	// fmt.Println(twoSum(nums, target))

	// nums = []int{3, 3}
	// target = 6
	// fmt.Println(twoSum(nums, target))

	nums := []int{1, 3, 4, 2}
	target := 6
	fmt.Println(twoSum(nums, target))
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

func twoSum(nums []int, target int) []int {
	elementMap := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		elementMap[nums[i]] = i
	}
	fmt.Println(elementMap)
	for index, val := range nums {
		if elementMap[target-val] != 0 && index != elementMap[target-val] {
			return []int{elementMap[target-val], index}
		}
	}

	return nil
}
