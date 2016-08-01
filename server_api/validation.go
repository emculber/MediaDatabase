package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

func (userKeys *UserKeys) validate() error {
	err := db.QueryRow(
		"SELECT user_keys.id, "+
			"registered_user.id, "+
			"registered_user.username, "+
			"role_permissions.id, "+
			"role_permissions.access, "+
			"role.id, "+
			"role.role, "+
			"permissions.id, "+
			"permissions.permission "+
			"FROM user_keys, "+
			"registered_user, "+
			"role_permissions, "+
			"role, "+
			"permissions "+
			"WHERE user_keys.user_id             = registered_user.id "+
			"AND user_keys.role_permissions_id   = role_permissions.id "+
			"AND role_permissions.role_id        = role.id "+
			"AND role_permissions.permissions_id = permissions.id "+
			"AND user_keys.key                   = $1", &userKeys.Key).Scan(&userKeys.Id, &userKeys.User.Id, &userKeys.User.Username, &userKeys.RolePermissions.Id, &userKeys.RolePermissions.access, &userKeys.RolePermissions.Role.Id, &userKeys.RolePermissions.Role.Role, &userKeys.RolePermissions.Permission.Id, &userKeys.RolePermissions.Permission.Permission)

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
