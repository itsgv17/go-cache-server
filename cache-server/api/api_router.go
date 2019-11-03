package api

import "net/http"

//Handler used for routing
func Handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/cache/v1/get", CacheGetHandleFunc)
	r.HandleFunc("/cache/v1/post", CachePostHandleFunc)
	r.HandleFunc("/cache/v1/delete", CacheDeleteHandleFunc)
	return r
}
