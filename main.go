/* Web application allowing users to convert URL to shortened version
Example: https:goo.gl/
*/

package main

import (
	"fmt"
	"hash/crc32"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Counter struct {
	Clicks  int
	FullURL string
}

func main() {

	store := map[string]*Counter{}
	formTemplate := template.Must(template.New("formTemplate").Parse(form))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			err := formTemplate.Execute(w, store)
			if err != nil {
				log.Println("Unable to open form.", err)
			}
		} else {
			shortCode := strings.TrimPrefix(r.URL.Path, "/")
			site, ok := store[shortCode]
			if ok {
				site.Clicks = site.Clicks + 1

				http.Redirect(w, r, site.FullURL, http.StatusFound)
			} else {
				http.NotFound(w, r)
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
			store[shortCode] = x
			w.Write([]byte(shortCode))

		}
	})

	http.ListenAndServe(":8080", nil)
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
		{{range $key, $value := $}}
		<p><a href="http://localhost:8080/{{$key}}"> {{$key}}</a> | {{$value.FullURL}} | {{$value.Clicks}}</p>
		{{end}}
	</div>
  </body>
</html>
`
