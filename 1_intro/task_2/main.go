package main

import (
	"fmt"
	"strings"
)

var punctuations = []string{"?", ".", ",", "'", ":", ";", "\""}

func GetWordCount(s string) map[string]int {
	wordCounter := map[string]int{}
	words := strings.Split(s, " ")
	for _, word := range words {
		word = strings.ToLower(word)
		for _, punct := range punctuations {
			word = strings.ReplaceAll(word, punct, "")
		}

		wordCounter[word]++
	}

	return wordCounter
}

func IsPalindrome(s string) bool {
	word := strings.TrimSpace(strings.ToLower(s))
	for _, punct := range punctuations {
		word = strings.ReplaceAll(word, punct, "")
	}

	i, j := 0, len(word)-1
	for i < j {
		if word[i] != word[j] {
			return false
		}
		i++
		j--
	}

	return true
}

func main() {
	f1t1 := "the red hare jumped over the sly fox"
	f1t2 := "Hello hello hello? he:llo"
	fmt.Printf("Expected len 7, got len %v\n", len(GetWordCount(f1t1)))
	fmt.Printf("Expected len 1, got len %v\n", len(GetWordCount(f1t2)))

	f2t1 := "this is not a palindrome   "
	f2t2 := " this siht   "
	fmt.Printf("Expected isPalindrome: false, got %v\n", IsPalindrome(f2t1))
	fmt.Printf("Expected isPalindrome: true, got %v\n", IsPalindrome(f2t2))
}
