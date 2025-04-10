package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Choice struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Choice `json:"options"`
}

func main() {
	handler, err := generateHandleFunc()

	if err != nil {
		log.Fatal("couldn't parse JSON")
	}

	http.Handle("/", handler)
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func generateHandleFunc() (http.HandlerFunc, error) {
	jsonData, fileReadErr := os.ReadFile("./gopher.json")

	if fileReadErr != nil {
		return nil, errors.New("couldn't read JSON file")
	}

	dec := json.NewDecoder(strings.NewReader(string(jsonData)))

	arcs := map[string]Arc{}

	for {
		if err := dec.Decode(&arcs); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {

		t := template.Must(template.ParseFiles("./template.html"))

		for key, value := range arcs {
			if r.URL.Path == fmt.Sprintf("/%s", key) {
				t.ExecuteTemplate(w, "Page", value)
				return
			}
		}

		http.NotFound(w, r)
	}, nil
}
