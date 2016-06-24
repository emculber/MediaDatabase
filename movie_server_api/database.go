package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/emculber/database_access/postgresql"
)

var db *sql.DB

type OmdbapiData struct {
	Id      int
	Imdb_id string
	Title   string
	Year    string
}

type User struct {
	Id       int
	User_key string
	Username string
}

type UsersMovie struct {
	Id           int
	User         User
	Omdbapi      OmdbapiData
	Movie_width  string
	Movie_height string
	Video_codac  string
	Audio_codac  string
	Container    string
	Frame_rate   string
	Aspect_ratio string
}

type UsersMovieJson struct {
	Username      string
	Movie_title   string
	Movie_imdb_id string
	Movie_width   string
	Movie_height  string
	Video_codac   string
	Audio_codac   string
	Container     string
	Frame_rate    string
	Aspect_ratio  string
}

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

func InsertNewMovie(omdbapi OmdbapiData) (OmdbapiData, error) {
	err := db.QueryRow(`insert into movie_list (imdb_id, movie_title, movie_year) values($1, $2, $3) returning id`, omdbapi.Imdb_id, omdbapi.Title, omdbapi.Year).Scan(&omdbapi.Id)
	if err != nil {
		return omdbapi, err
	}
	return omdbapi, nil
}

func InsertNewUsersMovie(users_movie UsersMovie) (UsersMovie, error) {
	err := db.QueryRow(`insert into users_movies (movie_list_id, user_id, width, height, video_codac, audio_codac, container, frame_rate, aspect_ratio) values($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`, users_movie.Omdbapi.Id, users_movie.User.Id, users_movie.Movie_width, users_movie.Movie_height, users_movie.Video_codac, users_movie.Audio_codac, users_movie.Container, users_movie.Frame_rate, users_movie.Aspect_ratio).Scan(&users_movie.Id)
	if err != nil {
		return users_movie, err
	}
	return users_movie, nil
}

func ReadUserMovies(user User) []UsersMovieJson {
	statement := fmt.Sprintf("SELECT registered_users.username, movie_list.movie_title, movie_list.imdb_id, users_movies.width, users_movies.height, users_movies.video_codac, users_movies.audio_codac, users_movies.container, users_movies.frame_rate, users_movies.aspect_ratio FROM users_movies, registered_users, movie_list WHERE user_id=%d AND registered_users.id = users_movies.user_id AND movie_list.id = users_movies.movie_list_id", user.Id)
	//TODO: Error Checking
	movies, _, _ := postgresql_access.QueryDatabase(db, statement)
	movies_list := []UsersMovieJson{}
	for _, movie := range movies {
		single_movie := UsersMovieJson{}
		single_movie.Username = movie[0].(string)
		single_movie.Movie_title = movie[1].(string)
		single_movie.Movie_imdb_id = movie[2].(string)
		single_movie.Movie_width = movie[3].(string)
		single_movie.Movie_height = movie[4].(string)
		single_movie.Video_codac = movie[5].(string)
		single_movie.Audio_codac = movie[6].(string)
		single_movie.Container = movie[7].(string)
		single_movie.Frame_rate = movie[8].(string)
		single_movie.Aspect_ratio = movie[9].(string)
		movies_list = append(movies_list, single_movie)

	}
	return movies_list
}
