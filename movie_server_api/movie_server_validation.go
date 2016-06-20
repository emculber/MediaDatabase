package main

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
)

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
