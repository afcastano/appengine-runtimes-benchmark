package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	if r.URL.Path == "/" {
		fmt.Fprintf(w, "Golang service running")
		return
	}

	if r.URL.Path != "/entities" {
		fmt.Fprintf(w, "Not found")
		return
	}

	ctx := appengine.NewContext(r)
	log.Infof(ctx, "Querying dummy entities")

	// if only one expected
	random2 := r.URL.Query().Get("random2")
	log.Infof(ctx, "Fetching where random 2 > %s", random2)

	random2val, _ := strconv.Atoi(random2)

	q := datastore.NewQuery("DummyEntity").Filter("random2 >=", random2val).Filter("random2 <", random2val+10000).Limit(10)
	var entities []DummyEntity

	if _, err := q.GetAll(ctx, &entities); err != nil {
		fmt.Fprintf(w, "Error %v", err)
		log.Errorf(ctx, "Error fetching entities: %v", err)
		return
	}

	log.Infof(ctx, "Found %d entities", len(entities))

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(entities)
}

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main() // Starts the server to receive requests
}
