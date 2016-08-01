package main

import (
	"encoding/json"
	"net/http"
)

func InitExternalSources() {
}

func (registeredMovie *RegisteredMovie) getMovieData() error {
	url := "http://www.omdbapi.com/?i=" + registeredMovie.Imdb_id
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&registeredMovie); err != nil {
		return err
	}

	if err = registeredMovie.OK(); err != nil {
		return err
	}
	return nil
}
