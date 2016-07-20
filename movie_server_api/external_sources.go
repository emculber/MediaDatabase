package main

import (
	"encoding/json"
	"net/http"
)

func InitExternalSources() {
}

func OmdbApi(imdb_id string) (RegisteredMovie, error) {
	RegisteredMovie := RegisteredMovie{}

	url := "http://www.omdbapi.com/?i=" + imdb_id
	r, err := http.Get(url)
	if err != nil {
		return RegisteredMovie, err
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&RegisteredMovie); err != nil {
		return RegisteredMovie, err
	}

	if err = RegisteredMovie.OK(); err != nil {
		return RegisteredMovie, err
	}
	return RegisteredMovie, nil
}
