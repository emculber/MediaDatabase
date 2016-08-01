package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/emculber/database_access/postgresql"
)

var db *sql.DB

var databaseSchema = []string{
	"CREATE TABLE registered_user(id serial primary key, username varchar)",
	"CREATE TABLE role(id serial primary key, role varchar)",
	"CREATE TABLE permissions(id serial primary key, permission varchar)",
	"CREATE TABLE role_permissions(id serial primary key, role_id integer references role(id), permissions_id integer references permissions(id), access varchar)",
	"CREATE TABLE user_keys(id serial primary key, user_id integer references registered_user(id), role_permissions_id integer references role_permissions(id), key varchar)",
	"CREATE TABLE registered_movie(id serial primary key, imdb_id varchar, title varchar, year varchar)",
	"CREATE TABLE movie_list(id serial primary key, registered_movie_id integer references registered_movie(id), user_id integer references registered_user(id), width varchar, height varchar, video_codac varchar, audio_codac varchar, container varchar, frame_rate varchar, aspect_ratio varchar)",
	"CREATE TABLE accepted_movie(id serial primary key, user_requested_id integer references registered_user(id), user_accepted_id integer references registered_user(id), registered_movie_id integer references registered_movie(id))",
	"CREATE TABLE exclude_movie(id serial primary key, user_id integer references registered_user(id), registered_movie_id integer references registered_movie(id))",
	"CREATE TABLE requested_movie(id serial primary key, user_id integer references registered_user(id), registered_movie_id integer references registered_movie(id))",
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

func (userRole *UserRole) getUserKey() error {
	err := db.QueryRow("select user_keys.key from user_keys, registered_user where user_keys.user_id = registered_user.id and registered_user.username = $1", userRole.User.Username).Scan(&userRole.Key)
	if err != nil {
		return err
	}
	return nil
}

func (registeredMovie *RegisteredMovie) RegisterNewMovie() error {
	err := db.QueryRow(`insert into registered_movie (imdb_id, title, year) values($1, $2, $3) returning id`, registeredMovie.Imdb_id, registeredMovie.Title, registeredMovie.Year).Scan(&registeredMovie.Id)
	if err != nil {
		return err
	}
	return nil
}

func (movie *MovieList) RegisterNewMovie() error {
	err := db.QueryRow(`insert into movie_list (registered_movie_id, user_id, width, height, video_codac, audio_codac, container, frame_rate, aspect_ratio) values($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`, movie.RegisteredMovie.Id, movie.UserRole.User.Id, movie.Movie_width, movie.Movie_height, movie.Video_codac, movie.Audio_codac, movie.Container, movie.Frame_rate, movie.Aspect_ratio).Scan(&movie.Id)
	if err != nil {
		return err
	}
	return nil
}

func (userRole *UserRole) ReadUserMovies() TransportMovies {
	statement := fmt.Sprintf("SELECT registered_movie.title, registered_movie.imdb_id, registered_movie.year, movie_list.width, movie_list.height, movie_list.video_codac, movie_list.audio_codac, movie_list.container, movie_list.frame_rate, movie_list.aspect_ratio FROM movie_list, user_keys, registered_movie WHERE movie_list.registered_movie_id = registered_movie.id AND user_keys.user_id = %d", userRole.User.Id)
	//TODO: Error Checking
	movies, _, _ := postgresql_access.QueryDatabase(db, statement)
	movies_list := TransportMovies{}
	movies_list.userRole = *userRole
	for _, movie := range movies {
		single_movie := MovieList{}
		single_movie.RegisteredMovie.Title = movie[0].(string)
		single_movie.RegisteredMovie.Imdb_id = movie[1].(string)
		single_movie.RegisteredMovie.Year = movie[2].(string)
		single_movie.Movie_width = movie[3].(string)
		single_movie.Movie_height = movie[4].(string)
		single_movie.Video_codac = movie[5].(string)
		single_movie.Audio_codac = movie[6].(string)
		single_movie.Container = movie[7].(string)
		single_movie.Frame_rate = movie[8].(string)
		single_movie.Aspect_ratio = movie[9].(string)
		movies_list.movieList = append(movies_list.movieList, single_movie)

	}
	return movies_list
}

func (userRole *UserRole) getAllMovies() TransportMovies {
	statement := fmt.Sprintf("SELECT registered_movie.title, registered_movie.imdb_id, registered_movie.year, movie_list.width, movie_list.height, movie_list.video_codac, movie_list.audio_codac, movie_list.container, movie_list.frame_rate, movie_list.aspect_ratio FROM movie_list, user_keys, registered_movie WHERE movie_list.registered_movie_id = registered_movie.id AND user_keys.user_id = %d", userRole.User.Id)
	//TODO: Error Checking
	movies, _, _ := postgresql_access.QueryDatabase(db, statement)
	movies_list := TransportMovies{}
	movies_list.userRole = *userRole
	for _, movie := range movies {
		single_movie := MovieList{}
		single_movie.RegisteredMovie.Title = movie[0].(string)
		single_movie.RegisteredMovie.Imdb_id = movie[1].(string)
		single_movie.RegisteredMovie.Year = movie[2].(string)
		single_movie.Movie_width = movie[3].(string)
		single_movie.Movie_height = movie[4].(string)
		single_movie.Video_codac = movie[5].(string)
		single_movie.Audio_codac = movie[6].(string)
		single_movie.Container = movie[7].(string)
		single_movie.Frame_rate = movie[8].(string)
		single_movie.Aspect_ratio = movie[9].(string)
		movies_list.movieList = append(movies_list.movieList, single_movie)

	}
	return movies_list
}