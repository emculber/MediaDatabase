package main

import "errors"

type RegisteredMovie struct {
	Id      int    `json:"-"`
	Imdb_id string `json:"imdbID"`
	Title   string `json:"Title"`
	Year    string `json:"Year"`
}

type AcceptedMovie struct {
	Id                int
	UserRequestedId   string
	UserAcceptedId    string
	RegisteredMovieId string
}

type ExludeMovie struct {
	Id                int
	UserId            string
	RegisteredMovieId string
}

type RequestedMovie struct {
	Id                int
	UserId            string
	RegisteredMovieId string
}

type User struct {
	Id       int
	Username string
}

type Role struct {
	Id   int
	Role string
}

type Permissions struct {
	Id         int
	Permission string
}

type RolePermissions struct {
	Id         int
	Role       Role
	Permission Permissions
	access     int
}

type UserKeys struct {
	Id              int
	User            User
	RolePermissions RolePermissions
	Key             string
}

type MovieList struct {
	Id              int
	UserKeys        UserKeys
	RegisteredMovie RegisteredMovie
	Movie_width     string
	Movie_height    string
	Video_codac     string
	Audio_codac     string
	Container       string
	Frame_rate      string
	Aspect_ratio    string
}

type TransportMovies struct {
	userKeys  UserKeys
	movieList []MovieList
}

func (movieList *MovieList) OK() error {
	if len(movieList.UserKeys.Key) == 0 {
		return errors.New("No Key")
	}
	if len(movieList.RegisteredMovie.Imdb_id) == 0 {
		return errors.New("No IMDB ID")
	}
	if len(movieList.Movie_width) == 0 {
		return errors.New("No Movie Width")
	}
	if len(movieList.Movie_height) == 0 {
		return errors.New("No Movie Height")
	}
	if len(movieList.Video_codac) == 0 {
		return errors.New("No Video Codac")
	}
	if len(movieList.Audio_codac) == 0 {
		return errors.New("No Audio Codac")
	}
	if len(movieList.Container) == 0 {
		return errors.New("No Container")
	}
	if len(movieList.Frame_rate) == 0 {
		return errors.New("No Frame Rate")
	}
	if len(movieList.Aspect_ratio) == 0 {
		return errors.New("No Aspect Ratio")
	}
	return nil
}

func (userKeys *UserKeys) OK() error {
	if len(userKeys.Key) == 0 {
		return errors.New("No Key")
	}
	return nil
}

func (registeredMovie *RegisteredMovie) OK() error {
	if len(registeredMovie.Imdb_id) == 0 {
		return errors.New("No IMDB ID")
	}
	if len(registeredMovie.Title) == 0 {
		return errors.New("No Title")
	}
	if len(registeredMovie.Year) == 0 {
		return errors.New("No Year")
	}
	return nil
}
