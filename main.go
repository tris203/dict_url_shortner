package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
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
	router := http.NewServeMux()
	staticFS := http.FileServer(http.Dir("./static"))
	router.HandleFunc("GET /", getRoot)
	router.HandleFunc("POST /shorten", getShorten)
	router.HandleFunc("POST /expand", getExpand)
	router.Handle("GET /static/", http.StripPrefix("/static/", staticFS))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := http.Server{
		Addr:    ":8080",
		Handler: Logging(router),
	}

	fmt.Printf("Listening on port %s\n", port)
	server.ListenAndServe()
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080/"
	}
	switch r.URL.Path {
	case "/about":
		getAbout(w)
	case "/":
		tmpl := template.Must(template.ParseFiles("templ/base.html", "templ/index.html", "templ/previous.html", "templ/footer.html"))
		tmpl.Execute(w, IndexData{BaseURL: baseURL, ExampleShortURL: "IXqwWqIXt", Title: "URL Shortener"})
	case "/favicon.ico":
		http.NotFound(w, r)
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
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	var expandedURL string
	expandedURL = expandUrl(url, words)

	data := ExpandData{ExpandURL: expandedURL, Title: "URL Redirect"}
	tmpl.Execute(w, data)
}
