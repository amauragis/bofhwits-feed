package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type FeedConfig struct {
	Mysql struct {
		Host string
		DB   string
		User string
		Pass string
	}
}

type BadPost struct {
	Username  string
	Time      string
	Requestor string
	Post      string
}

func getConfig() (FeedConfig, error) {
	path := "./feed.yaml"

	cfg := FeedConfig{}

	source, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Read file failure: %v\n", err)
		return cfg, err
	}

	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		log.Printf("Unmarshal YAML failure: %v\n", err)
		return cfg, err
	}

	return cfg, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, post []*BadPost) {
	t, err := template.ParseFiles("views/" + tmpl + ".html")
	if err != nil {
		log.Printf("Parse err: %v\n", err)
	}
	err = t.Execute(w, post)
	if err != nil {
		log.Printf("Execute err: %v\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	posts := []*BadPost{
		&BadPost{Username: "dumb", Time: "yesterday", Requestor: "me", Post: "lolcats"},
		&BadPost{Username: "buttes", Time: "today", Requestor: "you", Post: "Murder she wrote"},
		&BadPost{Username: "donges", Time: "tomorrow", Requestor: "her", Post: "idk"},
	}
	renderTemplate(w, "feed", posts)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
