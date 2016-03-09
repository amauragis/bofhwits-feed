package main

import (
	"database/sql"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/amauragis/sanitize"

	"gopkg.in/yaml.v2"
	// database driver tomfoolery
	_ "github.com/go-sql-driver/mysql"
)

type feedConfig struct {
	Mysql struct {
		Host string
		DB   string
		User string
		Pass string
	}
}

type badPost struct {
	Username  string
	Time      string
	Requestor string
	Post      string
}

var configs feedConfig

func getConfig(path string) error {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Read file failure: %v\n", err)
		return err
	}

	err = yaml.Unmarshal(source, &configs)
	if err != nil {
		log.Printf("Unmarshal YAML failure: %v\n", err)
		return err
	}

	return nil
}

func mysqlOpen() (*sql.DB, error) {
	return sql.Open("mysql",
		configs.Mysql.User+":"+configs.Mysql.Pass+
			"@tcp("+configs.Mysql.Host+":3306)/"+configs.Mysql.DB)
}

func queryPosts(sqlcon *sql.DB) (posts []*badPost, err error) {
	posts = make([]*badPost, 0)
	sqlQuery := `SELECT requestor, ts, user, post FROM bofhwits_posts ORDER BY ts DESC`
	stmt, err := sqlcon.Prepare(sqlQuery)
	if err != nil {
		log.Printf("Prepare Statement Failure: %v\n", err)
		return posts, err
	}
	defer stmt.Close()
	log.Printf("Querying: %s", sqlQuery)
	rows, err := stmt.Query()
	if err != nil {
		log.Printf("Query Failure: %v\n", err)
		return posts, err
	}

	var (
		requestor string
		ts        string
		user      string
		post      string
	)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&requestor, &ts, &user, &post)
		if err != nil {
			log.Printf("Row scan err: %v", err)
		}
		thisPost := badPost{
			Username:  sanitize.HTML(user),
			Time:      sanitize.HTML(ts),
			Requestor: sanitize.HTML(requestor),
			Post:      sanitize.HTML(post),
		}
		posts = append(posts, &thisPost)
		// log.Printf("Got row: | %s | %s | %s | %s |\n", requestor, ts, user, post)
	}
	err = rows.Err()
	if err != nil {
		log.Printf("rows Failure: %v\n", err)
		return posts, err
	}

	return posts, nil

}

func renderTemplate(w http.ResponseWriter, tmpl string, post []*badPost) {
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

	// this is stupid
	log.Printf("Recieved request: %v\n                    From: %v\n",
		r.RequestURI, r.RemoteAddr)

	// posts := []*badPost{
	// 	&badPost{Username: "dumb", Time: "yesterday", Requestor: "me", Post: "lolcats"},
	// 	&badPost{Username: "buttes", Time: "today", Requestor: "you", Post: "Murder she wrote"},
	// 	&badPost{Username: "donges", Time: "tomorrow", Requestor: "her", Post: "idk"},
	// }

	db, err := mysqlOpen()
	if err != nil {
		log.Printf("mysql open err: %v\n", err)
	}
	posts, err := queryPosts(db)
	if err != nil {
		log.Printf("query err: %v\n", err)
	}
	renderTemplate(w, "feed", posts)
}

func main() {

	configFile := flag.String("c", "config/feed.yaml",
		"The path to the configuration file to use")
	logFile := flag.String("l", "",
		"The path to the log file to use. If not provided, uses stdout")

	flag.Parse()

	if *logFile != "" {
		file, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			log.Fatalf("Could not open log file!  Error: %v\n", err)
		}
		defer file.Close()
		log.SetOutput(file)
	}

	err := getConfig(*configFile)
	if err != nil {
		log.Printf("Could not load configuration file!")
	}

	log.Println("Registering handlers...")
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", http.NotFound)

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
