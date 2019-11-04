package main

import (
	"log"
	"net/http"

	"./api"
)

func main() {

	go func() {
		api.CacheEvictionScheadular(1)
	}()

	if err := http.ListenAndServe(":8080", api.Handler()); err != nil {
		log.Fatal(err)
	}
}
