package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type entry struct {
	id        uuid.UUID
	createdAt time.Time
	value     interface{}
}

type storage struct {
	store map[string]entry
	mutex sync.Mutex
}

var db *storage

func setHandler(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	key := params.Get("key")
	value := params.Get("value")

	db.mutex.Lock()
	db.store[key] = entry{
		id:        uuid.New(),
		createdAt: time.Now(),
		value:     value,
	}
	db.mutex.Unlock()

	fmt.Fprintf(w, "ok")
}

func getHandler(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	key := params.Get("key")

	db.mutex.Lock()
	entry := db.store[key]
	db.mutex.Unlock()

	fmt.Fprintf(w, "%s", entry.value)
}

func main() {
	db = &storage{store: make(map[string]entry)}

	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}