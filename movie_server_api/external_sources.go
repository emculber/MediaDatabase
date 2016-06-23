package main

import (
	"encoding/json"
	"net/http"
)

type OmdbapiData struct {
	Imdb_id string
	Title   string
	Year    string
}

func InitExternalSources() {
}

func Omdbapi(imdb_id string) (OmdbapiData, error) {
	omdbapiData := OmdbapiData{}

	url := "http://www.omdbapi.com/?i=" + imdb_id
	r, err := http.Get(url)
	if err != nil {
		return omdbapiData, err
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&omdbapiData)
	return omdbapiData, nil
}
