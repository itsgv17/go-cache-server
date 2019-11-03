package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCacheGetRouter(t *testing.T) {

	server := httptest.NewServer(Handler())
	defer server.Close()

	testTable := []struct {
		name       string
		collection string
		key        string
		status     int
		err        string
	}{
		{name: "collection header missing", key: "key", status: http.StatusBadRequest, err: `{"err":"Collection request header not found"}`},
		{name: "key header missing", collection: "collection", status: http.StatusBadRequest, err: `{"err":"Key request header not found"}`},
		{name: "cache not found", collection: "collection", key: "noKey", status: http.StatusNotFound, err: fmt.Sprintf(`{"err":"Key %q not found in Collection %q"}`, "noKey", "collection")},
		{name: "cache found", collection: "collection", key: "key", status: http.StatusOK},
	}

	var innerMap = map[string]CachedObject{}

	cachedObject := CachedObject{CreatedAt: time.Now(), Value: `{"name":"gobi"}`}
	innerMap["key"] = cachedObject
	twoLevelMap["collection"] = innerMap

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/cache/v1/get", server.URL), nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			req.Header.Set("collection", testCase.collection)
			req.Header.Set("key", testCase.key)

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Request not successful: %v", err)
			}
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if testCase.err != "" {

				if res.StatusCode != testCase.status {
					t.Errorf("expected status %v; got %v", testCase.status, res.StatusCode)
				}
				if msg := string(bytes.TrimSpace(b)); msg != testCase.err {
					t.Errorf("expected message %q; got %q", testCase.err, msg)
				}
				return
			}

			if res.StatusCode != testCase.status {
				t.Errorf("expected status %v; got %v", testCase.status, res.Status)
			}

		})
	}

}

func TestCachePostRouter(t *testing.T) {

	server := httptest.NewServer(Handler())
	defer server.Close()

	testTable := []struct {
		name       string
		collection string
		key        string
		json       string
		status     int
		err        string
	}{
		{name: "collection header missing", key: "key", status: http.StatusBadRequest, err: `{"err":"Collection request header not found"}`},
		{name: "key header missing", collection: "collection", status: http.StatusBadRequest, err: `{"err":"Key request header not found"}`},
		{name: "Posting invalid Json", collection: "collection", key: "key", status: http.StatusUnsupportedMediaType, json: "name:goibinath", err: `{"err":"Invalid Json payload"}`},
		{name: "posting valid json", collection: "collection", key: "key", status: http.StatusOK, json: `{"name":"gobi"}`},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			var jsonbytes = []byte(testCase.json)

			req, err := http.NewRequest("POST", fmt.Sprintf("%s/cache/v1/post", server.URL), bytes.NewBuffer(jsonbytes))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			req.Header.Set("collection", testCase.collection)
			req.Header.Set("key", testCase.key)

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Request not successful: %v", err)
			}
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if testCase.err != "" {

				if res.StatusCode != testCase.status {
					t.Errorf("expected status %v; got %v", testCase.status, res.StatusCode)
				}
				if msg := string(bytes.TrimSpace(b)); msg != testCase.err {
					t.Errorf("expected message %q; got %q", testCase.err, msg)
				}
				return
			}

			if res.StatusCode != testCase.status {
				t.Errorf("expected status %v; got %v", testCase.status, res.Status)
			}

		})
	}

}

func TestCacheDeleteRouter(t *testing.T) {

	server := httptest.NewServer(Handler())
	defer server.Close()

	testTable := []struct {
		name       string
		collection string
		key        string
		status     int
		err        string
	}{
		{name: "collection header missing", key: "key", status: http.StatusBadRequest, err: `{"err":"Collection request header not found"}`},
		{name: "key header missing", collection: "collection", status: http.StatusBadRequest, err: `{"err":"Key request header not found"}`},
		{name: "cache not found", collection: "collection", key: "noKey", status: http.StatusNotFound, err: fmt.Sprintf(`{"err":"Key %q not found in Collection %q"}`, "noKey", "collection")},
		{name: "cache deleted", collection: "collection", key: "key", status: http.StatusOK},
	}

	var innerMap = map[string]CachedObject{}

	cachedObject := CachedObject{CreatedAt: time.Now(), Value: `{"name":"gobi"}`}
	innerMap["key"] = cachedObject
	twoLevelMap["collection"] = innerMap

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/cache/v1/delete", server.URL), nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			req.Header.Set("collection", testCase.collection)
			req.Header.Set("key", testCase.key)

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Request not successful: %v", err)
			}
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if testCase.err != "" {

				if res.StatusCode != testCase.status {
					t.Errorf("expected status %v; got %v", testCase.status, res.StatusCode)
				}
				if msg := string(bytes.TrimSpace(b)); msg != testCase.err {
					t.Errorf("expected message %q; got %q", testCase.err, msg)
				}
				return
			}

			if res.StatusCode != testCase.status {
				t.Errorf("expected status %v; got %v", testCase.status, res.Status)
			}

		})
	}
}
