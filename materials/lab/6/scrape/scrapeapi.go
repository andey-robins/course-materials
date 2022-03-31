package scrape

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"scrape/logging"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"
)

//==========================================================================\\

// Helper function walk function, modfied from Chap 7 BHG to enable passing in of
// additional parameter http responsewriter; also appends items to global Files and
// if responsewriter is passed, outputs to http

func walkFn(w http.ResponseWriter) filepath.WalkFunc {
	addedFileCount := 0

	return func(path string, f os.FileInfo, err error) error {
		w.Header().Set("Content-Type", "application/json")

		for _, r := range regexes {
			if r.MatchString(path) {
				var tfile FileInfo
				dir, filename := filepath.Split(path)
				tfile.Filename = string(filename)
				tfile.Location = string(dir)

				if !slices.Contains(Files, tfile) {
					Files = append(Files, tfile)
					addedFileCount++
				}

				if w != nil && len(Files) > 0 {
					w.Write([]byte(`"` + strconv.FormatInt(int64(addedFileCount), 10) + `":  `))
					json.NewEncoder(w).Encode(tfile)
					w.Write([]byte(`,`))
				}

				logging.Log("[+] HIT: "+path, 2)
			}
		}

		return nil
	}
}

func walkFn2(w http.ResponseWriter, query string) filepath.WalkFunc {
	r := regexp.MustCompile(query)
	addedFileCount := 0

	return func(path string, f os.FileInfo, err error) error {
		w.Header().Set("Content-Type", "application/json")

		if r.MatchString(path) {
			var tfile FileInfo
			dir, filename := filepath.Split(path)
			tfile.Filename = string(filename)
			tfile.Location = string(dir)

			if !slices.Contains(Files, tfile) {
				Files = append(Files, tfile)
				addedFileCount++
			}

			if w != nil && len(Files) > 0 {
				w.Write([]byte(`"` + strconv.FormatInt(int64(addedFileCount), 10) + `":  `))
				json.NewEncoder(w).Encode(tfile)
				w.Write([]byte(`,`))
			}

			logging.Log("[+] HIT: "+path, 2)
		}
		return nil

	}
}

//==========================================================================\\

func APISTATUS(w http.ResponseWriter, r *http.Request) {
	logging.Log("Entering"+r.URL.Path+" end point", 1)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status" : "API is up and running ", "ready": true,`))
	var regexstrings []string

	for _, regex := range regexes {
		regexstrings = append(regexstrings, regex.String())
	}

	w.Write([]byte(` "regexs" :`))
	json.NewEncoder(w).Encode(regexstrings)
	w.Write([]byte(`}`))
	log.Println(regexes)

}

func MainPage(w http.ResponseWriter, r *http.Request) {
	logging.Log("Entering"+r.URL.Path+" end point", 1)
	w.Header().Set("Content-Type", "text/html")

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, `<html>
        <body>
            <H1>Welcome to my awesome File page</H1>
            <p>
                This API can be used to examine a filesystem <br>
                The following endpoints are defined: <br>
                - /api-status -> check status. "ready": true means its good to go. <br>
                - /indexer -> Index files according to active regular expressions. 
								Optionally, pass in a query string with a regex to only search that one. <br>
                - /search -> Pass in a query string to search for a specific file name. <br>
                - /addsearch/{regex} -> Add the URL encoded regular expression to the searches. <br> 
                - /clear -> Clear all stored regex values. <br>
                - /reset -> Reset to default regex values.
            </p>
        </body>`)
}

func FindFile(w http.ResponseWriter, r *http.Request) {
	logging.Log("Entering"+r.URL.Path+" end point", 1)
	q, ok := r.URL.Query()["q"]

	w.WriteHeader(http.StatusOK)
	if ok && len(q[0]) > 0 {
		log.Printf("Entering search with query=%s", q[0])

		found := false
		for _, File := range Files {
			if File.Filename == q[0] {
				json.NewEncoder(w).Encode(File.Location)
				found = true
			}
		}

		if !found {
			w.Write([]byte("No matching files found. Please try again."))
		}

	} else {
		// didn't pass in a search term, show all that you've found
		w.Write([]byte(`"files":`))
		json.NewEncoder(w).Encode(Files)
	}
}

func IndexFiles(w http.ResponseWriter, r *http.Request) {
	logging.Log("Entering"+r.URL.Path+" end point", 1)
	w.Header().Set("Content-Type", "application/json")

	location, locOK := r.URL.Query()["location"]
	regexStr, regexOK := r.URL.Query()["regex"]
	rootDir := "/Users" // macOS safe starting point. User /home for ubuntu and a different OS for windows :p

	if locOK && len(location[0]) > 0 {
		w.WriteHeader(http.StatusOK)

	} else {
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(`{ "parameters" : {"required": "location",`))
		w.Write([]byte(`"optional": "regex"},`))
		w.Write([]byte(`"examples" : { "required": "/indexer?location=/xyz",`))
		w.Write([]byte(`"optional": "/indexer?location=/xyz&regex=(i?).md"}}`))
		return
	}

	//wrapper to make "nice json"
	w.Write([]byte(`{ `))

	// if regexOK
	if regexOK {
		if err := filepath.Walk(rootDir+location[0], walkFn2(w, `(i?)`+regexStr[0])); err != nil {
			log.Panicln(err)
		}
	} else if err := filepath.Walk(rootDir+location[0], walkFn(w)); err != nil {
		log.Panicln(err)
	}

	//wrapper to make "nice json"
	w.Write([]byte(` "status": "completed"} `))

}

func AddSearch(w http.ResponseWriter, r *http.Request) {
	logging.Log("Entering"+r.URL.Path+" end point", 1)
	w.Header().Set("Content-Type", "application/json")
	regex := mux.Vars(r)["regex"]
	if regex != "" {
		addRegEx("(?i)" + regex)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func Clear(w http.ResponseWriter, r *http.Request) {
	logging.Log("Entering"+r.URL.Path+" end point", 1)
	w.Header().Set("Content-Type", "application/json")
	clearRegEx()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func Reset(w http.ResponseWriter, r *http.Request) {
	logging.Log("Entering"+r.URL.Path+" end point", 1)
	w.Header().Set("Content-Type", "application/json")
	resetRegEx()
	Files = make([]FileInfo, 0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
