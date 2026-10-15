//go:build !solution

package main

import (
	"fmt"
	"flag"
	"sync"
	
	"hash/fnv"
	"encoding/json"
	"net/http"
)


type State struct {
	mutex sync.RWMutex
	urlMap map[string]string	
}

type RequestBody struct {
	Url string `json:"url"`
}
type ResponseBody struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

func generateKey(url string) string {
	h := fnv.New32a()
	h.Write([]byte(url))
	return fmt.Sprintf("%d", h.Sum32())
}

func serve(state *State, port int) {
	shortenHandler := func(w http.ResponseWriter, r *http.Request){
		var body RequestBody
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20) 

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
	        return
		}
		key := generateKey(body.Url)

		state.mutex.RLock()
		_, ok := state.urlMap[key]
		state.mutex.RUnlock()

		if !ok {
			state.mutex.Lock()
    		state.urlMap[key] = body.Url
    		state.mutex.Unlock()
		}
		

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&ResponseBody{Url: body.Url, Key: key})
	}

	goHandler := func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

        state.mutex.RLock()
		url, ok := state.urlMap[key]
		state.mutex.RUnlock()
		
		if !ok {
			http.Error(w, "key not found", http.StatusNotFound)
		} else {
			http.Redirect(w, r, url, http.StatusFound)
		}
	}
	
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", shortenHandler)
	mux.HandleFunc("GET /go/{key}", goHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)

}

func main() {
	port := flag.Int("port", 8080, "Port to listen")
	flag.Parse()

	serve(&State{urlMap: make(map[string]string)}, *port)
}
