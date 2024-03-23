package main

import (
	"testing"
)

func TestBase62Encode(t *testing.T) {
	var tests = []struct {
		input    int
		expected string
	}{
		{0, "000"},
		{1, "001"},
		{10, "00a"},
		{61, "00Z"},
		{62, "010"},
		{63, "011"},
		{100, "01C"},
		{1000, "0g8"},
		{10000, "2Bi"},
		{100000, "q0U"},
		{238327, "ZZZ"},
	}

	for _, test := range tests {
		result := base62Encode(test.input)
		if result != test.expected {
			t.Errorf("base62Encode(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func BenchmarkBase62Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		base62Encode(238327)
	}
}

func TestBase62Decode(t *testing.T) {
	var tests = []struct {
		input    string
		expected int
	}{
		{"000", 0},
		{"001", 1},
		{"00a", 10},
		{"00Z", 61},
		{"010", 62},
		{"011", 63},
		{"01C", 100},
		{"0g8", 1000},
		{"2Bi", 10000},
		{"q0U", 100000},
		{"ZZZ", 238327},
	}

	for _, test := range tests {
		result := base62Decode(test.input)
		if result != test.expected {
			t.Errorf("base62Decode(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func BenchmarkBase62Decode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		base62Decode("ZZZ")
	}
}

func TestGetDictionary(t *testing.T) {
	dictionary := getDictionary()
	if len(dictionary) == 0 {
		t.Fatalf("getDictionary() = %v; want a non-empty dictionary", dictionary)
	}
}

func BenchmarkGetDictionary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getDictionary()
	}
}

func TestShortenURL(t *testing.T) {
	var words = getDictionary()
	url := "https://www.google.com"
	shortURL := shortenURL(url, words)
	if shortURL == "" {
		t.Errorf("ShortenURL(%s) = %s; want a valid short URL", url, shortURL)
	}
	if shortURL == url {
		t.Errorf("ShortenURL(%s) = %s; was not shortened", url, shortURL)
	}
	//check the shortURL is valid -  IXrIM2LyYLz3LyWIXt
	if shortURL != "IXrIXsIM2LyYLz3LyWIXt" {
		t.Errorf("ShortenURL(%s) = %s; did not match expected value", url, shortURL)
	}

	url2 := "https://twitter.com"
	shortURL2 := shortenURL(url2, words)
	if shortURL2 != "IXrwWqIXt" {
		t.Errorf("ShortenURL(%s) = %s; did not match expected value", url2, shortURL2)
	}

	if shortURL == shortURL2 {
		t.Errorf("ShortenURL(%s) = %s; was not unique", url, shortURL)
	}
}

func BenchmarkShortenURL(b *testing.B) {
	var words = getDictionary()
	for i := 0; i < b.N; i++ {
		shortenURL("https://www.google.com", words)
	}
}

func TestExpandURL(t *testing.T) {
	var words = getDictionary()
	url := "https://www.google.com"
	shortURL := shortenURL(url, words)
	expandedURL := expandUrl(shortURL, words)
	if expandedURL != url {
		t.Errorf("ExpandURL(%s) = %s; want %s", shortURL, expandedURL, url)
	}

	url2 := "https://twitter.com"
	shortURL2 := shortenURL(url2, words)
	expandedURL2 := expandUrl(shortURL2, words)
	if expandedURL2 != url2 {
		t.Errorf("ExpandURL(%s) = %s; want %s", shortURL2, expandedURL2, url2)
	}

	if expandedURL == expandedURL2 {
		t.Errorf("ExpandURL(%s) = %s; was not unique", shortURL, expandedURL)
	}

}

func TestExpandEmptyURL(t *testing.T) {
	var words = getDictionary()
	expandedURL := expandUrl("", words)
	if expandedURL != "" {
		t.Errorf("ExpandURL(%s) = %s; want %s", "", expandedURL, "")
	}
}

func BenchmarkExpandURL(b *testing.B) {
	var words = getDictionary()
	for i := 0; i < b.N; i++ {
		expandUrl("IXrIXsIM2LyYLz3LyWIXt", words)
	}
}
