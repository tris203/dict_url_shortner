package main

import (
	"strings"
)

func isValidURL(url string) bool {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	return true
}

func getBaseURL(baseURL string) string {
	if baseURL == "" {
		baseURL = "http://localhost:8080/"
	}
	return baseURL
}

func padSlice(slice []string, length int) []string {
	for i := len(slice); i < length; i++ {
		slice = append(slice, "YOUWILLNEVERFINDMEINAURLORSTUFFAOIJDOJOIAWJDOIJDWIJDWIJ")
	}
	return slice
}
