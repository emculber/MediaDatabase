package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

/*
TODO:
Get User Key,
Get User list,
Get User Movie List
*/

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func getUserKey(w http.ResponseWriter, r *http.Request) {

	userKeys := UserKeys{}

	r.ParseForm()
	userKeys.Key = r.PostFormValue("key")
	userKeys.User.Username = r.PostFormValue("username")

	if err := userKeys.getUserKey(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Getting User Key")
	}

	if err := json.NewEncoder(w).Encode(userKeys.User); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Invalid Request!")
	}
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()

	userKeys := UserKeys{}
	userKeys.Key = r.PostFormValue("key")

	if err := userKeys.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Validating User")
	}

	movies := userKeys.ReadUserMovies()

	if err := json.NewEncoder(w).Encode(movies.movieList); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Invalid Request!")
	}
}

func getAllRegesteredMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()

	userKeys := UserKeys{}
	userKeys.Key = r.PostFormValue("key")
	userKeys.User.Username = r.PostFormValue("username")

	if err := userKeys.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Validating User")
	}

	movies := userKeys.getAllMovies()

	if err := json.NewEncoder(w).Encode(movies.movieList); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Invalid Request!")
	}
}

func addMovieToUserMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()

	userKeys := UserKeys{}
	userKeys.Key = r.PostFormValue("key")

	registeredMovie := RegisteredMovie{}
	registeredMovie.Imdb_id = strings.ToLower(r.PostFormValue("imdb_id"))

	movieList := MovieList{}
	movieList.UserKeys = userKeys
	movieList.RegisteredMovie = registeredMovie
	movieList.Movie_width = r.PostFormValue("movie_width")
	movieList.Movie_height = r.PostFormValue("movie_height")
	movieList.Video_codac = r.PostFormValue("video_codac")
	movieList.Audio_codac = r.PostFormValue("audio_codac")
	movieList.Container = r.PostFormValue("container")
	movieList.Frame_rate = r.PostFormValue("frame_rate")
	movieList.Aspect_ratio = r.PostFormValue("aspect_ratio")

	log.WithFields(log.Fields{
		"Movies": movieList,
	}).Info("Movie List")

	if err := movieList.OK(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Movie List OK")
	}

	if err := movieList.UserKeys.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Validating User")
	}

	if err := movieList.UserKeys.RolePermissions.checkAccess("write"); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error")
	}

	if err := movieList.RegisteredMovie.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error")
		if err := movieList.RegisteredMovie.getMovieData(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error")
		} else {
			log.WithFields(log.Fields{
				"Movie": movieList.RegisteredMovie,
			}).Info("Registered Movie")
		}
		if err := movieList.RegisteredMovie.RegisterNewMovie(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error")
		} else {
			log.WithFields(log.Fields{
				"Movie": movieList.RegisteredMovie,
			}).Info("Registered Movie")
		}
	} else {
		log.WithFields(log.Fields{
			"Movie": movieList.RegisteredMovie,
		}).Info("Movie Already Registered")
	}

	if err := movieList.validate(); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error")
		if err := movieList.RegisterNewMovie(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error")
		} else {
			log.WithFields(log.Fields{
				"Movies": movieList,
			}).Info("Movie List")
		}
	}
	w.Write([]byte("OK"))
}
