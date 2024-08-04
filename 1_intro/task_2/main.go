package main

import "strings"

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
	}

	return true
}

func main() {}
