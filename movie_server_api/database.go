package main

import (
	"database/sql"
	"os"

	"github.com/emculber/database_access/postgresql"
)

var db *sql.DB

func InitDatabase() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	db, err = postgresql_access.ConfigFilePathAutoConnect(dir + "/config.json")
	if err != nil {
		panic(err)
	}
}

func InsertNewMovie(omdbapi OmdbapiData) (int, error) {
	var id int
	err := db.QueryRow(`insert into movie_list (imdb_id, movie_title, movie_year) values($1, $2, $3) returning id`, omdbapi.Imdb_id, omdbapi.Title, omdbapi.Year).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}
