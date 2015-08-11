package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type dataStore interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
}

type keyGenerator interface {
	New() (string, error)
}

type Config struct {
	RootPath     string
	Port         string
	VCAPServices string
}

type Application struct {
	RootPath     string
	Port         string
	DataStore    dataStore
	KeyGenerator keyGenerator
	GridHandler  http.Handler
}

func getRedisCloudConfig(vcap_services string) (string, string, error) {
	var services map[string][]struct {
		Name        string
		Credentials struct {
			Hostname string
			Password string
			Port     string
		}
	}
	err := json.Unmarshal([]byte(vcap_services), &services)
	if err != nil {
		return "", "", err
	}

	creds := services["rediscloud"][0].Credentials
	return fmt.Sprintf("%s:%s", creds.Hostname, creds.Port), creds.Password, nil
}

func NewApplication(config Config) (*Application, error) {
	gridHTML, err := ioutil.ReadFile(filepath.Join(config.RootPath, "webclient", "grid.html"))
	if err != nil {
		return nil, err
	}
	redisHost, redisPassword, err := getRedisCloudConfig(config.VCAPServices)
	if err != nil {
		panic(err)
		return nil, err
	}
	dataStore := &RedisDataStore{Host: redisHost, Password: redisPassword}
	err = dataStore.Connect()
	if err != nil {
		return nil, err
	}

	keyGenerator := KeyGenerator{dataStore}
	gridHandler := GridHandler{gridHTML, dataStore, keyGenerator}
	return &Application{
		RootPath:     config.RootPath,
		Port:         config.Port,
		DataStore:    dataStore,
		KeyGenerator: keyGenerator,
		GridHandler:  gridHandler,
	}, nil
}

func (a *Application) Boot() error {
	mux := http.NewServeMux()
	mux.Handle("/grid/assets/", a.getAssetsHandler())
	mux.Handle("/grid/", a.GridHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		a.homePageHandler(w, req)
	})

	return http.ListenAndServe(":"+a.Port, mux)
}

func (a *Application) getAssetsHandler() http.Handler {
	return http.StripPrefix("/grid/assets/",
		http.FileServer(http.Dir(a.RootPath+"/webclient/assets")))
}

func (a *Application) homePageHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, `<html><body><a href="/grid/new">New Grid</a></body></html>`)
}
