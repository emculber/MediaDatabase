package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func checkUserKey(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	key := r.PostFormValue("key")

	user := getUsername(key)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Invalid Request!")
	}
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user := UserRole{}
	user.Key = r.PostFormValue("user_key")

	/*
		var isValidated bool
		var err error
		isValidated, user, err = validateUserId(user)
		if err != nil || !isValidated {
			log.WithFields(log.Fields{
				"Logger Fields":    api_logger_fields,
				"Error":            http.StatusBadRequest,
				"Validation Error": err,
				"User":             user,
				"Registered User":  isValidated,
			}).Error("Invalid User")
			http.Error(w, "Invalid Request -> Invalid User", http.StatusBadRequest)
			return
		}

		movies := ReadUserMovies(user)

		if err = json.NewEncoder(w).Encode(movies); err != nil {
			log.WithFields(log.Fields{
				"Logger Fields": api_logger_fields,
				"Error":         err,
			}).Error("Invalid Request!")
			panic(err)
		}
	*/
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

	fmt.Println(permissions["movie"])

	/*
		var isValidated bool
		var err error
		isValidated, users_movie.User, err = validateUserId(users_movie.User)
		if err != nil || !isValidated {
			log.WithFields(log.Fields{
				"Logger Fields":    api_logger_fields,
				"Error":            http.StatusBadRequest,
				"Validation Error": err,
				"User":             users_movie.User,
				"Data":             users_movie,
				"Registered User":  isValidated,
			}).Error("Invalid User")
			http.Error(w, "Invalid Request -> Invalid User", http.StatusBadRequest)
			return
		}

		isValidated, users_movie.Omdbapi, err = validateImdbId(users_movie.Omdbapi)
		if err != nil || !isValidated {
			users_movie.Omdbapi, err = Omdbapi(users_movie.Omdbapi.Imdb_id)
			if err != nil {
				log.WithFields(log.Fields{
					"Logger Fields": api_logger_fields,
					"Error":         err,
					"Movie":         users_movie.Omdbapi,
					"Data":          users_movie,
				}).Error("Error Finding Movie")
			}
			users_movie.Omdbapi, err = InsertNewMovie(users_movie.Omdbapi)
			if err != nil {
				log.WithFields(log.Fields{
					"Logger Fields": api_logger_fields,
					"Error":         err,
					"Movie":         users_movie.Omdbapi,
					"Data":          users_movie,
				}).Error("Error Inserting Movie")
			}
			log.WithFields(log.Fields{
				"Validation Error": err,
				"Movie":            users_movie.Omdbapi,
				"Data":             users_movie,
				"Registered Movie": isValidated,
			}).Info("Added movie")
		}

		log.WithFields(log.Fields{
			"Data": users_movie,
		}).Info("Adding Users Movie")

		isValidated, err = validateUsersMovie(users_movie)
		if isValidated {
			log.WithFields(log.Fields{
				"Logger Fields": api_logger_fields,
				"Data":          users_movie,
			}).Info("Movie Already Added")
		} else {
			users_movie, err = InsertNewUsersMovie(users_movie)
			if err != nil {
				log.WithFields(log.Fields{
					"Logger Fields":  api_logger_fields,
					"Error":          http.StatusBadRequest,
					"Database Error": err,
					"Data":           users_movie,
				}).Error("Error adding data to database")
				http.Error(w, "Invalid Request -> Error Adding Movie", http.StatusBadRequest)
				return
			}
			log.WithFields(log.Fields{
				"Data": users_movie,
			}).Info("Movie Added")
		}
	*/
	w.Write([]byte("OK"))
}
