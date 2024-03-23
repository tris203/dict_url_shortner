package main

import (
	"testing"
)

func TestIsValidURL(t *testing.T) {
	url := "http://example.com"
	if !isValidURL(url) {
		t.Fatalf("isValidURL(%s) = should validate", url)
	}
}

func TestIsValidURLInvalid(t *testing.T) {
	url := "example.com"
	if isValidURL(url) {
		t.Fatalf("isValidURL(%s) = should not validate", url)
	}
}

func TestGetBaseURL(t *testing.T) {
	baseURL := "http://example.com"
	if getBaseURL(baseURL) != baseURL {
		t.Fatalf("getBaseURL(%s) = %s, want %s", baseURL, getBaseURL(baseURL), baseURL)
	}
}

func TestNullBaseURL(t *testing.T) {
	baseURL := ""
	expected := "http://localhost:8080/"
	if getBaseURL(baseURL) != expected {
		t.Fatalf("getBaseURL(%s) = %s, want %s", baseURL, getBaseURL(baseURL), expected)
	}
}
