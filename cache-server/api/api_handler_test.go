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

func TestCacheGetHandleFunc(t *testing.T) {
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
			req, err := http.NewRequest("GET", "localhost:8080/cache/v1/get", nil)
			req.Header.Set("collection", testCase.collection)
			req.Header.Set("key", testCase.key)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			CacheGetHandleFunc(rec, req)

			res := rec.Result()
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

func TestCachePostHandleFunc(t *testing.T) {
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

			req, err := http.NewRequest("POST", "localhost:8080/cache/v1/post", bytes.NewBuffer(jsonbytes))
			req.Header.Set("collection", testCase.collection)
			req.Header.Set("key", testCase.key)

			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			CachePostHandleFunc(rec, req)

			res := rec.Result()
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

func TestCacheDeleteHandleFunc(t *testing.T) {
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
			req, err := http.NewRequest("DELETE", "localhost:8080/cache/v1/delete", nil)
			req.Header.Set("collection", testCase.collection)
			req.Header.Set("key", testCase.key)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			CacheDeleteHandleFunc(rec, req)

			res := rec.Result()
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
