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
	"CREATE TABLE users_roles(id serial primary key, user_id integer references registered_user(id), role_id integer references role(id), key varchar)",
	"CREATE TABLE permissions(id serial primary key, permission varchar)",
	"CREATE TABLE role_permissions(id serial primary key, role_id integer references role(id), permissions_id integer references permissions(id), access varchar)",
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

func getUsername(key string) User {
	statement := fmt.Sprintf("select registered_user.username from registered_user, user_roles where user_roles.key='%s'", key)
	username, _, _ := postgresql_access.QueryDatabase(db, statement)
	user := User{}
	user.Username = username[0][0].(string)

	return user
}

/*
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
*/
