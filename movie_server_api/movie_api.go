package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/emculber/database_access/postgresql"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

var config Configuration
var db *sql.DB

type databaseInfo struct {
	Host     string
	Port     int
	Username string
	Password string
	Dbname   string
}

type Configuration struct {
	Db databaseInfo
}

type user struct {
	id       int
	username string
	key      string
}

type new_movie_list struct {
	Title string
	Year  string
}

func init() {

	path := os.Getenv("GOPATH")
	var filePath string = path + "/logs/api_logs/media_database.log"

	log.SetFormatter(&log.TextFormatter{})
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	log_file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.WithFields(log.Fields{
			"File Path": filePath,
			"Error":     err.Error(),
		}).Error("Error Opening File")
	}

	log.SetOutput(log_file)

	log.WithFields(log.Fields{
		"Log Format": "Text Format",
		"Log level":  "Info",
		"Log Output": log_file,
	}).Info("Format, Level, Output set")

	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s/%s", config.Db.Username, config.Db.Password, config.Db.Host, config.Db.Dbname)
	fmt.Println(dbUrl)
	log.WithFields(log.Fields{
		"Postgresql Server": dbUrl,
	}).Info("Connecting to Postgresql Server")
	db, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.WithFields(log.Fields{
			"Postgresql Server": dbUrl,
			"Error":             err,
		}).Error("ERROR -> Connecting to Postgresql Server")
		panic(err)
	}
}

func validateUserId(key string) (string, error) {
	var userInfo user
	err := db.QueryRow("select id, username, key from registered_users where key = $1", key).Scan(&userInfo.id, &userInfo.username, &userInfo.key)
	switch {
	case err != nil:
		log.WithFields(log.Fields{
			"key":   key,
			"Error": err,
		}).Error("ERROR -> Validating User Id")
		return "", err
	}

	log.WithFields(log.Fields{
		"key":      userInfo.key,
		"username": userInfo.username,
		"id":       userInfo.id,
	}).Info("User Accessed")

	return strconv.Itoa(userInfo.id), nil
}

func validateImdbId(id string) (string, error) {
	var dbId int
	err := db.QueryRow("select id from movie_list where imdb_id = $1", id).Scan(&dbId)
	switch {
	case err != nil:
		log.WithFields(log.Fields{
			"IMDB Id": id,
			"Error":   err,
		}).Error("ERROR -> Validating ImdbId Id")
		return "", err
	}

	return strconv.Itoa(dbId), nil
}

func addMovieToList(imdb_id string) {
	url := "http://www.omdbapi.com/?i=" + imdb_id
	new_movie := new_movie_list{}
	r, err := http.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&new_movie)
	log.WithFields(log.Fields{
		"json": new_movie,
	}).Info("Test")
	var id int
	err = db.QueryRow(`insert into movie_list (imdb_id, movie_title, movie_year) values($1, $2, $3) returning id`, imdb_id, new_movie.Title, new_movie.Year).Scan(&id)
	if err != nil {
		log.WithFields(log.Fields{
			"json":    new_movie,
			"imdb id": id,
		}).Error("Error adding data to database")
		return
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")
	log.WithFields(log.Fields{
		"_Method":     r.Method,
		"_IP address": ip[0],
		"_Port":       ip[1],
		"IP:Port":     r.RemoteAddr,
	}).Info("Test was hit")
	w.Write([]byte("OK"))
}

func addMovieToUserMovies(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")
	if r.Method != "POST" {
		log.WithFields(log.Fields{
			"Method":     r.Method,
			"IP address": ip[0],
			"Port":       ip[1],
			"Error":      http.StatusMethodNotAllowed,
		}).Error("Invalid Request!")
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
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
	pixel_format := r.PostFormValue("pixel_format")
	aspect_ratio := r.PostFormValue("aspect_ratio")

	if user_key == "" || imdb_id == "" || movie_width == "" ||
		movie_height == "" || audio_codac == "" || video_codac == "" ||
		container == "" || frame_rate == "" || pixel_format == "" || aspect_ratio == "" {

		log.WithFields(log.Fields{
			"_Method":      r.Method,
			"_IP address":  ip[0],
			"_Port":        ip[1],
			"_Error":       http.StatusBadRequest,
			"user_key":     user_key,
			"imdb_id":      imdb_id,
			"movie_width":  movie_width,
			"movie_height": movie_height,
			"audio_codac":  audio_codac,
			"video_codac":  video_codac,
			"container":    container,
			"frame_rate":   frame_rate,
			"pixel_format": pixel_format,
			"aspect_ratio": aspect_ratio,
		}).Error("Empty Content")
		http.Error(w, "Invalid Request!", http.StatusBadRequest)
		return
	}

	user_id, err := validateUserId(user_key)
	if err != nil {
		log.WithFields(log.Fields{
			"_Method":           r.Method,
			"_IP address":       ip[0],
			"_Port":             ip[1],
			"_Error":            http.StatusBadRequest,
			"_Validation Error": err,
			"user_key":          user_key,
		}).Error("Invalid User")
		http.Error(w, "Invalid Request! Invalid User!", http.StatusBadRequest)
		return
	}

	movie_list_id, err := validateImdbId(imdb_id)
	if err != nil {
		log.WithFields(log.Fields{
			"_Method":           r.Method,
			"_IP address":       ip[0],
			"_Port":             ip[1],
			"_Validation Error": err,
			"imdb_id":           imdb_id,
		}).Info("Adding movie")
		addMovieToList(imdb_id)
		return
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
		"pixel_format":       pixel_format,
		"aspect_ratio":       aspect_ratio,
	}).Info("Adding Movie")

	//Check for movie if all ready added

	var id int
	err = db.QueryRow(`insert into users_movies (movie_list_id, user_id, width, height, video_codac, audio_codac, container, frame_rate, pixel_format, aspect_ratio) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`, movie_list_id, user_id, movie_width, movie_height, video_codac, audio_codac, container, frame_rate, pixel_format, aspect_ratio).Scan(&id)
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

	w.Write([]byte("OK"))
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")
	if r.Method != "POST" {
		log.WithFields(log.Fields{
			"Method":     r.Method,
			"IP address": ip[0],
			"Port":       ip[1],
			"Error":      http.StatusMethodNotAllowed,
		}).Error("Invalid Request!")
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}

	statement := fmt.Sprintf("select id, movie_list_id, user_id, width, height, video_codac, audio_codac, container, frame_rate, pixel_format, aspect_ratio from users_movies")

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
	fmt.Println(config)
	http.HandleFunc("/api/getallmovies", getAllMovies)
	http.HandleFunc("/api/addmovie", addMovieToUserMovies)
	http.HandleFunc("/api/test", test)
	http.ListenAndServe(":8080", nil)
}
