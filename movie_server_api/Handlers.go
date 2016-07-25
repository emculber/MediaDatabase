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

	userRole := UserRole{}

	r.ParseForm()
	userRole.key = r.PostFormValue("key")
	userRole.User.Username = r.PostFormValue("username")

	user := getUsername(key)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Invalid Request!")
	}
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()

	userRole := UserRole{}
	userRole.Key = r.PostFormValue("key")

	if err := userRole.validate(); err != nil {
		fmt.Println(err)
	}

	movies := userRole.ReadUserMovies()

	if err := json.NewEncoder(w).Encode(movies.movieList); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Invalid Request!")
	}
}

func getAllRegesteredMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()

	userRole := UserRole{}
	userRole.Key = r.PostFormValue("key")
	userRole.User.Username = r.PostFormValue("username")

	if err := userRole.validate(); err != nil {
		fmt.Println(err)
	}

	movies := getAllMovies()

	if err := json.NewEncoder(w).Encode(movies.movieList); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Invalid Request!")
	}
}

func addMovieToUserMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()

	userRole := UserRole{}
	userRole.Key = r.PostFormValue("key")

	registeredMovie := RegisteredMovie{}
	registeredMovie.Imdb_id = strings.ToLower(r.PostFormValue("imdb_id"))

	movieList := MovieList{}
	movieList.UserRole = userRole
	movieList.RegisteredMovie = registeredMovie
	movieList.Movie_width = r.PostFormValue("movie_width")
	movieList.Movie_height = r.PostFormValue("movie_height")
	movieList.Video_codac = r.PostFormValue("video_codac")
	movieList.Audio_codac = r.PostFormValue("audio_codac")
	movieList.Container = r.PostFormValue("container")
	movieList.Frame_rate = r.PostFormValue("frame_rate")
	movieList.Aspect_ratio = r.PostFormValue("aspect_ratio")

	fmt.Println(movieList)

	if err := movieList.OK(); err != nil {
		fmt.Println(err)
	}

	if err := movieList.UserRole.validate(); err != nil {
		fmt.Println(err)
	}

	if err := movieList.UserRole.RolePermissions.checkAccess("write"); err != nil {
		fmt.Println(err)
	}

	if err := movieList.RegisteredMovie.validate(); err != nil {
		fmt.Println(err)
		if err := movieList.RegisteredMovie.getMovieData(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(movieList.RegisteredMovie)
		}
		if err := movieList.RegisteredMovie.RegisterNewMovie(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Movie Registered")
			fmt.Println(movieList.RegisteredMovie)
		}
	} else {
		fmt.Println("Movie Already Registered")
	}

	if err := movieList.validate(); err != nil {
		fmt.Println(err)
		if err := movieList.RegisterNewMovie(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(movieList)
		}
	}
	w.Write([]byte("OK"))
}
