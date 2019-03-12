package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/appengine" // Required external App Engine library
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type DummyEntity struct {
	Id      string `datastore:"id"`
	Random1 string `datastore:"random1"`
	Random2 int    `datastore:"random2"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// if statement redirects all invalid URLs to the root homepage.
	// Ex: if URL is http://[YOUR_PROJECT_ID].appspot.com/FOO, it will be
	// redirected to http://[YOUR_PROJECT_ID].appspot.com.
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	ctx := appengine.NewContext(r)
	log.Infof(ctx, "Querying dummy entities")
	q := datastore.NewQuery("DummyEntity").Filter("random2 >", 10000).Filter("random2 <", 100000).Limit(10)
	var entities []DummyEntity

	// if _, err := q.GetAll(ctx, &entities); err != nil {
	// 	log.Errorf(ctx, "Error fetching entities: %v", err)
	// 	return
	// }

	for t := q.Run(ctx); ; {
		var entity DummyEntity
		key, err := t.Next(&entity)

		if err == datastore.Done {
			break
		}

		if err != nil {
			fmt.Fprintf(w, "Key=%v\nWidget=%#v\n\n", key, err)
			log.Errorf(ctx, "Error fetching entities: %v", err)
			return
		}

		entities = append(entities, entity)
		// fmt.Fprintf(w, "Key=%v\nWidget=%#v\n\n", key, entity)
	}
	log.Infof(ctx, "Found %d entities", len(entities))

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(entities)
}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main() // Starts the server to receive requests
}
