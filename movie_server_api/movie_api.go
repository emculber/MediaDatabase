package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/emculber/database_access/postgresql"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

type user struct {
	id       int
	username string
	key      string
}

func init() {

	InitLogger()
	InitDatabase()
	InitExternalSources()

}

func addMovieToList(imdb_id string) {
	omdbapiData, err := Omdbapi(imdb_id)
	id, err := InsertNewMovie(omdbapiData)
	fmt.Println(id, err)
}

func test(w http.ResponseWriter, r *http.Request) {

	api_logger_fields := ApiLoggerFields{}
	api_logger_fields.ip_address = r.RemoteAddr
	api_logger_fields.method_type = r.Method

	log.WithFields(log.Fields{
		"Logger Fields": loggerFields,
	}).Info("Test was hit")

	w.Write([]byte("OK"))
}

func addMovieToUserMovies(w http.ResponseWriter, r *http.Request) {

	api_logger_fields := ApiLoggerFields{}
	api_logger_fields.ip_address = r.RemoteAddr
	api_logger_fields.method_type = r.Method

	if r.Method != "POST" {
		log.WithFields(log.Fields{
			"Logger Fields": loggerFields,
			"Error":         http.StatusMethodNotAllowed,
		}).Error("Invalid Request!")
		http.Error(w, "Invalid Request -> Incorrect Method Call", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()

	user_key := r.PostFormValue("user_key")
	imdb_id := strings.ToLower(r.PostFormValue("imdb_id"))
	movie_width := r.PostFormValue("movie_width")
	movie_height := r.PostFormValue("movie_height")
	video_codac := r.PostFormValue("video_codac")
	audio_codac := r.PostFormValue("audio_codac")
	container := r.PostFormValue("container")
	frame_rate := r.PostFormValue("frame_rate")
	aspect_ratio := r.PostFormValue("aspect_ratio")

	if user_key == "" || imdb_id == "" || movie_width == "" ||
		movie_height == "" || audio_codac == "" || video_codac == "" ||
		container == "" || frame_rate == "" || aspect_ratio == "" {

		log.WithFields(log.Fields{
			"Logger Fields": loggerFields,
			"Error":         "addMovieToUserList -> Empty Value Detected",
			"user_key":      user_key,
			"imdb_id":       imdb_id,
			"movie_width":   movie_width,
			"movie_height":  movie_height,
			"audio_codac":   audio_codac,
			"video_codac":   video_codac,
			"container":     container,
			"frame_rate":    frame_rate,
			"aspect_ratio":  aspect_ratio,
		}).Error("Empty Content")
		http.Error(w, "Invalid Request -> Empty Content Was Detected", http.StatusBadRequest)
		return
	}

	user_id, err := validateUserId(user_key)
	if err != nil {
		log.WithFields(log.Fields{
			"Logger Fields":    loggerFields,
			"Error":            http.StatusBadRequest,
			"Validation Error": err,
			"user_key":         user_key,
		}).Error("Invalid User")
		http.Error(w, "Invalid Request -> Invalid User", http.StatusBadRequest)
		return
	}

	movie_list_id, err := validateImdbId(imdb_id)
	if err != nil {
		omdbapiData, err := Omdbapi(imdb_id)
		id, err := InsertNewMovie(omdbapiData)
		log.WithFields(log.Fields{
			"Logger Fields":    loggerFields,
			"Validation Error": err,
			"imdb_id":          imdb_id,
		}).Info("Adding movie")
		addMovieToList(imdb_id)
	}

	log.WithFields(log.Fields{
		"_Method":            r.Method,
		"_IP address":        ip[0],
		"_Port":              ip[1],
		"_Error":             http.StatusBadRequest,
		"user_key":           user_key,
		"user_id":            user_id,
		"imdb_id":            imdb_id,
		"imdb_movie_list_id": movie_list_id,
		"movie_width":        movie_width,
		"movie_height":       movie_height,
		"audio_codac":        audio_codac,
		"video_codac":        video_codac,
		"container":          container,
		"frame_rate":         frame_rate,
		"aspect_ratio":       aspect_ratio,
	}).Info("Adding Movie")

	//Check for movie if all ready added

	isAdded_statment := fmt.Sprintf("select movie_list_id, user_id from users_movies, movie_list where users_movies.movie_list_id = movie_list.id and users_movies.movie_list_id=%s and users_movies.user_id = %s", movie_list_id, user_id)

	isAdded, count, _ := postgresql_access.QueryDatabase(db, isAdded_statment)
	if count != 0 {
		if isAdded[0][0] == movie_list_id && isAdded[0][1] == user_id {
			log.WithFields(log.Fields{
				"_Method":            r.Method,
				"_IP address":        ip[0],
				"_Port":              ip[1],
				"_Error":             http.StatusBadRequest,
				"user_key":           user_key,
				"user_id":            user_id,
				"imdb_id":            imdb_id,
				"imdb_movie_list_id": movie_list_id,
				"movie_width":        movie_width,
				"movie_height":       movie_height,
				"audio_codac":        audio_codac,
				"video_codac":        video_codac,
				"container":          container,
				"frame_rate":         frame_rate,
				"aspect_ratio":       aspect_ratio,
			}).Info("Movie Already Added")
		} else {
			var id int
			err = db.QueryRow(`insert into users_movies (movie_list_id, user_id, width, height, video_codac, audio_codac, container, frame_rate, aspect_ratio) values($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`, movie_list_id, user_id, movie_width, movie_height, video_codac, audio_codac, container, frame_rate, aspect_ratio).Scan(&id)
			if err != nil {
				log.WithFields(log.Fields{
					"_Method":         r.Method,
					"_IP address":     ip[0],
					"_Port":           ip[1],
					"_Error":          http.StatusBadRequest,
					"_Database Error": err,
				}).Error("Error adding data to database")
				http.Error(w, "Invalid Request!", http.StatusBadRequest)
				return
			}
			log.WithFields(log.Fields{
				"id": id,
			}).Info("Movie Added")
		}
	} else {

		var id int
		err = db.QueryRow(`insert into users_movies (movie_list_id, user_id, width, height, video_codac, audio_codac, container, frame_rate, aspect_ratio) values($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`, movie_list_id, user_id, movie_width, movie_height, video_codac, audio_codac, container, frame_rate, aspect_ratio).Scan(&id)
		if err != nil {
			log.WithFields(log.Fields{
				"_Method":         r.Method,
				"_IP address":     ip[0],
				"_Port":           ip[1],
				"_Error":          http.StatusBadRequest,
				"_Database Error": err,
			}).Error("Error adding data to database")
			http.Error(w, "Invalid Request!", http.StatusBadRequest)
			return
		}
		log.WithFields(log.Fields{
			"id": id,
		}).Info("Movie Added")
	}
	w.Write([]byte("OK"))
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")
	if r.Method != "POST" {
		log.WithFields(log.Fields{
			"Method (Expected -> POST)": r.Method,
			"IP address":                ip[0],
			"Port":                      ip[1],
			"Error":                     http.StatusMethodNotAllowed,
		}).Error("Invalid Request!")
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	statement := fmt.Sprintf("select id, movie_list_id, user_id, width, height, video_codac, audio_codac, container, frame_rate, aspect_ratio from users_movies")

	movies, _, _ := postgresql_access.QueryDatabase(db, statement)

	if err := json.NewEncoder(w).Encode(movies); err != nil {
		log.WithFields(log.Fields{
			"Method":     r.Method,
			"IP address": ip[0],
			"Port":       ip[1],
			"Error":      err,
		}).Error("Invalid Request!")
		panic(err)
	}
}

func main() {
	http.HandleFunc("/api/getallmovies", getAllMovies)
	http.HandleFunc("/api/addmovie", addMovieToUserMovies)
	http.HandleFunc("/api/test", test)
	http.ListenAndServe(":8080", nil)
}
