package application

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type GridHandler struct {
	GridHTML     []byte
	DataStore    dataStore
	KeyGenerator keyGenerator
}

func getAbsURL(relativePath string, req *http.Request) string {
	newURL, _ := url.Parse(relativePath)
	return req.URL.ResolveReference(newURL).String()
}

func (g GridHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/grid/" {
		http.Redirect(w, req, getAbsURL("/", req), http.StatusMovedPermanently)
		return
	}
	if req.URL.Path == "/grid/new" {
		newKey, err := g.KeyGenerator.New()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		relativePath := "/grid/" + newKey
		http.Redirect(w, req, getAbsURL(relativePath, req), http.StatusTemporaryRedirect)
		return
	}
	if strings.HasSuffix(req.URL.Path, "/data") {
		g.DataHandler(w, req)
		return
	}
	if strings.HasSuffix(req.URL.Path, "/clone") {
		g.CloneHandler(w, req)
		return
	}

	key := strings.TrimPrefix(req.URL.Path, "/grid/")
	if _, err := g.DataStore.Get(key); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(g.GridHTML)
}

type GridElement struct {
	Layout struct {
		Col   int
		Row   int
		SizeX int
		SizeY int
	}
	Content struct {
		URL string
	}
}

type GridData []GridElement

func (g GridHandler) CloneHandler(w http.ResponseWriter, req *http.Request) {
	key := strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/grid/"), "/clone")
	val, err := g.DataStore.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	newKey, err := g.KeyGenerator.New()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = g.DataStore.Set(newKey, val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, req, getAbsURL("/grid/"+newKey, req), http.StatusTemporaryRedirect)
}

func (g GridHandler) DataHandler(w http.ResponseWriter, req *http.Request) {
	key := strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/grid/"), "/data")

	if req.Method == "GET" {
		val, err := g.DataStore.Get(key)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(val)
		return
	}

	if req.Method == "PUT" {
		value, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		var gridData []GridElement
		if err := json.Unmarshal(value, &gridData); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if err := g.DataStore.Set(key, value); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
