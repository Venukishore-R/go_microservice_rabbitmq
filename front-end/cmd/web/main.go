package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port 1600")
	err := http.ListenAndServe(":1600", nil)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}
}

//go:embed templates

var templateFS embed.FS

func render(w http.ResponseWriter, t string) {

	partials := []string{
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datas := struct {
		BrokerURL string
	}{
		BrokerURL: os.Getenv("BROKER_URL"),
		// BrokerURL: "http://10.106.8.149:1700",
	}

	if err := tmpl.Execute(w, datas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
