package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var words = getDictionary()

type ShortData struct {
	ShortURL        string
	OriginalLength  int
	ShortenedLength int
	BaseURL         string
}

type ExpandData struct {
	Title     string
	ExpandURL string
}

type IndexData struct {
	Title           string
	BaseURL         string
	ExampleShortURL string
}

func main() {
	staticFS := http.FileServer(http.Dir("./static"))
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/shorten", getShorten)
	http.HandleFunc("/expand", getExpand)
	http.Handle("/static/", http.StripPrefix("/static/", staticFS))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080/"
	}
	switch r.URL.Path {
	case "/about":
		getAbout(w)
		return
	case "/":
		tmpl := template.Must(template.ParseFiles("templ/base.html", "templ/index.html", "templ/previous.html", "templ/footer.html"))
		tmpl.Execute(w, IndexData{BaseURL: baseURL, ExampleShortURL: "IXqwWqIXt", Title: "URL Shortener"})
	default:
		getExpand(w, r)
	}
}

func getAbout(w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles("templ/base.html", "templ/about.html", "templ/footer.html"))
	tmpl.Execute(w, IndexData{Title: "About"})
}

func getShorten(w http.ResponseWriter, r *http.Request) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080/"
	}
	//get url from post body
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusUnprocessableEntity)
		return
	}
	if !isValidURL(url) {
		http.Error(w, "Invalid URL", http.StatusUnprocessableEntity)
		return
	}
	var shortURL string
	tmpl := template.Must(template.ParseFiles("templ/shorten.html", "templ/previous.html"))
	shortURL = shortenURL(url, words)
	data := ShortData{ShortURL: shortURL, OriginalLength: len(url), ShortenedLength: len(shortURL), BaseURL: baseURL}
	tmpl.Execute(w, data)

}

func getExpand(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templ/base.html", "templ/expand.html", "templ/footer.html"))
	url := r.URL.Path[1:]
	if url == "" {
		fmt.Println("URL is required")
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	var expandedURL string
	expandedURL = expandUrl(url, words)

	data := ExpandData{ExpandURL: expandedURL, Title: "URL Redirect"}
	tmpl.Execute(w, data)
}

func isValidURL(url string) bool {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	return true
}
