package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var cachedObjectCount int = 0

func getExpiryTimeInHrs() time.Duration {

	expiryTime := os.Getenv("EXPIRY_TIME")

	duration, err := time.ParseDuration(expiryTime)

	if err != nil {
		return 2
	}

	return duration

}

func getCacheLimit() int {

	maxCacheLimit := os.Getenv("MAX_CACHE")

	maxCache, err := strconv.Atoi(maxCacheLimit)

	if err != nil {
		return 1 << 10
	}

	return maxCache
}

// CachedObject type with cacheValue and createdAt
type CachedObject struct {
	CreatedAt time.Time   `json:"createdAt"`
	Value     interface{} `json:"value"`
}

// HTTPReturn type with Collection, Key and Value
type HTTPReturn struct {
	Collection string      `json:"collection"`
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
}

var (
	//twoLevelMap cache dataStructure
	twoLevelMap = map[string]map[string]CachedObject{}
	mutex       sync.Mutex
)

// IsExpired is used for checking cachedObject Expiration
func (c CachedObject) IsExpired() bool {

	elapsed := time.Now().Sub(c.CreatedAt.Add(time.Hour * getExpiryTimeInHrs()))

	if elapsed > 0.0 {
		return true
	}

	return false
}

// GetValue is used to return cached value
func (c CachedObject) GetValue() interface{} {

	return c.Value
}

// CacheGetHandleFunc is used to fetch the cached object
func CacheGetHandleFunc(w http.ResponseWriter, r *http.Request) {

	switch method := r.Method; method {

	case http.MethodGet:
		collection := r.Header.Get("collection")

		if strings.Compare(collection, "") == 0 {
			http.Error(w, fmt.Sprintf(`{"err":"Collection request header not found"}`), http.StatusBadRequest)
			return
		}

		key := r.Header.Get("key")

		if strings.Compare(key, "") == 0 {
			http.Error(w, fmt.Sprintf(`{"err":"Key request header not found"}`), http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		value, ok := twoLevelMap[collection][key]

		isExpired := value.IsExpired()
		if isExpired {
			delete(twoLevelMap[collection], key)
			cachedObjectCount--
		}

		if !ok || isExpired {
			http.Error(w, fmt.Sprintf(`{"err":"Key %q not found in Collection %q"}`, key, collection), http.StatusNotFound)
			return
		}

		httpReturn := HTTPReturn{
			Collection: collection,
			Key:        key,
			Value:      value.Value,
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(&httpReturn); err != nil {
			log.Printf("Unable to convert as Json %+v - %s", httpReturn, err)
		}

	default:
		http.Error(w, fmt.Sprintf(`{"err":"%q HTTP Method Not Supported. Supported HTTP Method(s) - %q"}`, method, http.MethodGet), http.StatusMethodNotAllowed)

	}
}

// CachePostHandleFunc is used to post the cache object
func CachePostHandleFunc(w http.ResponseWriter, r *http.Request) {

	switch method := r.Method; method {

	case http.MethodPost:

		collection := r.Header.Get("collection")

		if strings.Compare(collection, "") == 0 {
			http.Error(w, fmt.Sprintf(`{"err":"Collection request header not found"}`), http.StatusBadRequest)
			return
		}

		key := r.Header.Get("key")

		if strings.Compare(key, "") == 0 {
			http.Error(w, fmt.Sprintf(`{"err":"Key request header not found"}`), http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		dec := json.NewDecoder(r.Body)

		var jsonValue interface{}

		if err := dec.Decode(&jsonValue); err != nil {
			http.Error(w, `{"err":"Invalid Json payload"}`, http.StatusUnsupportedMediaType)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		if cachedObjectCount >= getCacheLimit() {
			http.Error(w, `{"err":"Cache Overflow"}`, http.StatusForbidden)
			return
		}

		cachedObject := CachedObject{CreatedAt: time.Now(), Value: jsonValue}
		innerMap, ok := twoLevelMap[collection]

		if !ok {

			twoLevelMap[collection] = map[string]CachedObject{}
			twoLevelMap[collection][key] = cachedObject
			cachedObjectCount++

		} else {

			_, ok = innerMap[key]
			if ok {
				innerMap[key] = cachedObject

			} else {
				innerMap[key] = cachedObject
				cachedObjectCount++
			}

		}

		httpReturn := HTTPReturn{
			Collection: collection,
			Key:        key,
			Value:      jsonValue,
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(httpReturn); err != nil {
			log.Printf("Unable to convert as Json %+v - %s", httpReturn, err)
		}
	default:
		http.Error(w, fmt.Sprintf(`{"err":"%q HTTP Method Not Supported. Supported HTTP Method(s) - %q"}`, method, http.MethodPost), http.StatusMethodNotAllowed)

	}

}

// CacheDeleteHandleFunc is used to remove cached object
func CacheDeleteHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {

	case http.MethodDelete:
		collection := r.Header.Get("collection")

		if strings.Compare(collection, "") == 0 {
			http.Error(w, fmt.Sprintf(`{"err":"Collection request header not found"}`), http.StatusBadRequest)
			return
		}

		key := r.Header.Get("key")

		if strings.Compare(key, "") == 0 {
			http.Error(w, fmt.Sprintf(`{"err":"Key request header not found"}`), http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		value, ok := twoLevelMap[collection][key]

		if !ok {
			http.Error(w, fmt.Sprintf(`{"err":"Key %q not found in Collection %q"}`, key, collection), http.StatusNotFound)
			return
		}

		delete(twoLevelMap[collection], key)
		cachedObjectCount--

		httpReturn := HTTPReturn{
			Collection: collection,
			Key:        key,
			Value:      value.Value,
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		if err := enc.Encode(&httpReturn); err != nil {
			log.Printf("Unable to convert as Json %+v - %s", httpReturn, err)
		}

	default:
		http.Error(w, fmt.Sprintf(`{"err":"%q HTTP Method Not Supported. Supported HTTP Method(s) - %q"}`, method, http.MethodDelete), http.StatusMethodNotAllowed)

	}
}

// CacheEvictionScheadular is used to evict expired cache
func CacheEvictionScheadular(interval time.Duration) {

	for {

		for collection := range twoLevelMap {

			for key := range twoLevelMap[collection] {

				mutex.Lock()

				if twoLevelMap[collection][key].IsExpired() {

					delete(twoLevelMap[collection], key)

					log.Printf("Key %q has been removed from collection %q by go routine ", key, collection)
				}

				mutex.Unlock()
			}
		}

		time.Sleep(time.Minute * interval)
	}
}
