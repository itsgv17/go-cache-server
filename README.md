# go_cache-server


Used GoLang and Nested Map (with 2 levels) uesed for implementation.
 ```
 
Outer Map:

var twoLevelMap = map[string]map[string]CachedObject{}
Cache Obj. Structure:
type CachedObject struct {
	CreatedAt time.Time   `json:"createdAt"`
	Value     interface{} `json:"value"`
}
Outer Map has keys as String and values as map[string]CachedObject{}

InnerMap( result value of Level 1 Map)-> contains actual cache key and cached object.

e.g.innerMap := twoLevelMap.get("CollectionKey");


Two Level Map data Structure example:

{
  "Registration": {
    "/{id=1}/register": {
      "name": "gobinath",
      "pass": "_#845HJHj"
    },
    "/{id=2}/register": {
      "name": "rama",
      "pass": "&*yyFH^8"
    }
  },
  "Payment": {
    "/{id=2}/pay": {
      "fromId": 123,
      "toId": 235,
      "amount": 100
    },
    "/{id=2}/credit": {
      "id": 123,
      "amount": 100
    },
    "/{id=3}/debit": {
      "id": 1345,
      "amount": 200
    }
  }
}

InOrder to fetch cached data from CacheServer, client has to send collection key as well as actual key.

Above approach will partition cache data based on collectionKey,
so the search will only happen on innerMap Map.
```

Implemented Unit Tests for Handler layer as well as Routing (Rest Api) Layer

Created three rest endpoints for Cache operation.
  - Fetching cached data (GET)
  ```
  
  package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func main() {

	url := "http://localhost:8080/cache/v1/get"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("collection", "collection")
	req.Header.Add("key", "key")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}

  ```
  - Posting data to cache server (POST)
  ```
  package main

import (
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
)

func main() {

	url := "http://localhost:8080/cache/v1/post"

	payload := strings.NewReader("{\n    \"serialNumber\": \"1234567890qwerty\",\n    \"latitude\": 123.001,\n    \"longitude\": \"123.002\"\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("collection", "collection")
	req.Header.Add("key", "key")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Length", "95")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}

  ```
  
  
  - Deleting cached data from Cache server (DELETE)
  ```
  package main

package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func main() {

	url := "http://localhost:8080/cache/v1/delete"

	req, _ := http.NewRequest("DELETE", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("collection", "collection")
	req.Header.Add("key", "key")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "localhost:8080")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Length", "0")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
  ```
  
  - Cache Eviction Scheadular
  ```
  
  It runs every one minute, removes expired cache objects
  
  Implemented approach iterates entire data structure & evicts expired cached object.
 
 
  ```
  
  
  
  
  
