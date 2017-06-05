package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

// Open a connection to Redis
var redisStorage = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_URL"),
	Password: "", // no password set
	DB:       0,  // use default DB
})

// -------- STRUCTURES

// Link to be stored
type Link struct {
	ID  int64
	URL string `json:"url"`
}

// Links array
type Links []Link

// RedirectHandler blabla
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	params := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(params) != 2 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	urlID, e := strconv.ParseInt(params[1], 0, 64)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url, err := redisStorage.Get(fmt.Sprintf("URL_%d", urlID)).Result()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	link := Link{
		ID:  urlID,
		URL: url,
	}
	linkj, _ := json.Marshal(link)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", linkj)
}

// RedirectionHandler blabla
func RedirectionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	params := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(params) != 2 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	urlID, e := strconv.ParseInt(params[1], 0, 64)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url, err := redisStorage.Get(fmt.Sprintf("URL_%d", urlID)).Result()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}

// ShorthenHandler blabla
func ShorthenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var link Link
	// Populate the user data
	json.NewDecoder(r.Body).Decode(&link)

	u, err := url.ParseRequestURI(link.URL)
	if err != nil && len(u.RawPath) > 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	key := getNextKey()
	// SET { URL_<key>: url } in Redis
	err2 := redisStorage.Set(fmt.Sprintf("URL_%d", key), link.URL, 0).Err()
	if err2 != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Add an Id
	link.ID = key
	// Marshal provided interface into JSON structure
	linkj, _ := json.Marshal(link)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", linkj)

}

func getNextKey() int64 {
	// INCR in Redis
	err := redisStorage.Incr("key").Err()
	if err != nil {
		panic(err)
	}

	// GET the key we just INCR'd
	key, _ := redisStorage.Get("key").Int64()

	return key
}
