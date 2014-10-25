package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	DBfile string
	port   string

	basePath  = os.Getenv("HOME") + "/.local/share/datacan/"
	htmlPath  = basePath + "htmls/"
	indexPage = htmlPath + "index.html"
	cssPath   = basePath + "styles/"
)

func init() {
	flag.StringVar(&DBfile, "db", basePath+"musicks.db", "The SQLite3 db file to use.")
	flag.StringVar(&port, "port", "19999", "Port to use for the service.")
	flag.Parse()

	db, err := sql.Open("sqlite3", DBfile)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS musicks (
            id INTEGER NOT NULL PRIMARY KEY,
            url TEXT,
            artist TEXT,
            title TEXT
        )`)
	if err != nil {
		panic(err)
	}
}

func writeDBHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", DBfile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	isBlank := func(s string) string {
		if s == "" {
			return "blank"
		}
		return s
	}
	lenLimit := func(s string) string {
		if len(s) > 100 {
			return "blank"
		}
		return s
	}
	url := lenLimit(isBlank(r.FormValue("url")))
	artist := lenLimit(isBlank(r.FormValue("artist")))
	title := lenLimit(isBlank(r.FormValue("title")))

	stmt, err := db.Prepare(`INSERT INTO musicks(url, artist, title) VALUES(?, ?, ?)`)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(url, artist, title)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
func readDBHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", DBfile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT url, artist, title FROM musicks`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var url, artist, title string
		rows.Scan(&url, &artist, &title)
		fmt.Fprintf(w, "%s\n%s\n%s\n\n", url, artist, title)
	}
	rows.Close()
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	cnt, err := ioutil.ReadFile(indexPage)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(cnt))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", writeDBHandler)
	http.HandleFunc("/view", readDBHandler)
	http.Handle(
		"/styles/",
		http.StripPrefix(
			"/styles/",
			http.FileServer(http.Dir(cssPath)),
		),
	)
	http.ListenAndServe(":"+port, nil)
}
