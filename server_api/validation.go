package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

func (userKeys *UserKeys) validate() error {
	err := userKeys.getUserInfo()
	fmt.Println(userKeys)

	if err != nil {
		log.WithFields(log.Fields{
			"User Role": userKeys,
			"Error":     err,
		}).Error("ERROR -> Validating User Id")
		return err
	}

	log.WithFields(log.Fields{
		"User Role": userKeys,
	}).Info("User Accessed")

	return nil
}

func (registeredMovie *RegisteredMovie) validate() error {
	err := db.QueryRow("select id, imdb_id, title, year from registered_movie where imdb_id = $1", registeredMovie.Imdb_id).Scan(&registeredMovie.Id, &registeredMovie.Imdb_id, &registeredMovie.Title, &registeredMovie.Year)
	if err != nil {
		log.WithFields(log.Fields{
			"Movie": registeredMovie,
			"Error": err,
		}).Error("ERROR -> Validating ImdbId Id")
		return err
	}

	return nil
}

func (movieList *MovieList) validate() error {
	err := db.QueryRow("select movie_list.id from movie_list, registered_user, registered_movie where registered_movie.id = movie_list.registered_movie_id and registered_user.id = movie_list.user_id and movie_list.registered_movie_id=$1 and movie_list.user_id=$2", movieList.RegisteredMovie.Id, movieList.UserKeys.User.Id).Scan(&movieList.Id)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("ERROR -> Validating ImdbId Id")
		return err
	}

	return nil
}
