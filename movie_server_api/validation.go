package main

import log "github.com/Sirupsen/logrus"

func validateUserId(userRole UserRole) (bool, User, error) {
	var user = User{}
	err := db.QueryRow("select id, username, key from registered_users where key = $1", userRole.Key).Scan(&user.Id, &user.Username)
	if err != nil {
		log.WithFields(log.Fields{
			"User":  user,
			"Error": err,
		}).Error("ERROR -> Validating User Id")
		return false, user, err
	}

	log.WithFields(log.Fields{
		"User": user,
	}).Info("User Accessed")

	return true, user, nil
}

func validateImdbId(registeredMovie RegisteredMovie) (bool, RegisteredMovie, error) {
	err := db.QueryRow("select id, imdb_id, movie_title, movie_year from movie_list where imdb_id = $1", registeredMovie.Imdb_id).Scan(&registeredMovie.Id, registeredMovie.Imdb_id, registeredMovie.Title, registeredMovie.Year)
	if err != nil {
		log.WithFields(log.Fields{
			"Movie": registeredMovie,
			"Error": err,
		}).Error("ERROR -> Validating ImdbId Id")
		return false, registeredMovie, err
	}

	return true, registeredMovie, nil
}

/*
func validateUsersMovie(users_movie UsersMovie) (bool, error) {
	validate_users_movie_statment := fmt.Sprintf("select movie_list_id, user_id from users_movies, movie_list where users_movies.movie_list_id = movie_list.id and users_movies.movie_list_id=%s and users_movies.user_id = %s", users_movie.Omdbapi.Id, users_movie.User.Id)
	isAdded, count, err := postgresql_access.QueryDatabase(db, validate_users_movie_statment)
	if count != 0 || err != nil {
		if isAdded[0][0] == users_movie.Omdbapi.Id && isAdded[0][1] == users_movie.User.Id {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, err
	}
	return false, err
}
*/
