package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func main() {
	file, err := os.Open("./gopher.json")
	if err != nil {
		panic(err)
	}
	fileDecoder := json.NewDecoder(file)
	var story map[string]Arc

	if err := fileDecoder.Decode(&story); err != nil {
		fmt.Println(err)
	}

	tmpl := template.Must(template.ParseFiles("layout.html"))
	arcHandler := ArcHandler(story, tmpl)

	fmt.Println("Statring the server on :80")
	http.ListenAndServe(":80", arcHandler)
}

func ArcHandler(story map[string]Arc, tmpl *template.Template) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if arc, ok := story[path[1:]]; ok {
			tmpl.Execute(rw, arc)
			return
		}
		tmpl.Execute(rw, story["intro"])
	}
}

type Option struct {
	Text    string `json:"text"`
	NextArc string `json:"arc"`
}

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}
