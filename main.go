/* Web application allowing users to convert URL to shortened version
Example: https:goo.gl/
*/

package main

import (
	"database/sql"
	"fmt"
	"hash/crc32"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {

	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	formTemplate := template.Must(template.New("formTemplate").Parse(form))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			counters, err := getCounters(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = formTemplate.Execute(w, counters)
			if err != nil {
				log.Println("Unable to open form.", err)
			}
		} else {
			shortCode := strings.TrimPrefix(r.URL.Path, "/")
			site, err := getCounter(db, shortCode)
			switch {
			case err == sql.ErrNoRows:
				http.NotFound(w, r)
			case err != nil:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			default:
				if err := incrementClicks(db, site); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				http.Redirect(w, r, site.FullURL, http.StatusFound)
			}
		}
	})

	crc32q := crc32.MakeTable(0xD5828281)
	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			shortenedURL := r.FormValue("url")

			x := &Counter{FullURL: shortenedURL}
			checkSum := crc32.Checksum([]byte(shortenedURL), crc32q)
			shortCode := fmt.Sprintf("%08x", checkSum)
			x.ShortURL = shortCode
			err := createCounter(db, x)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Write([]byte(shortCode))

		}
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

const form = `
<!DOCTYPE html>
<html>
  <body>
    <form method="post" action="/shorten">
      <input type="text" name="url"/>
      <button type="submit">Shorten</button>
    </form>
	<div>
		{{range $c := $}}
		<p><a href="/{{$c.ShortURL}}"> {{$c.ShortURL}}</a> | {{$c.FullURL}} | {{$c.Clicks}}</p>
		{{end}}
	</div>
  </body>
</html>
`
