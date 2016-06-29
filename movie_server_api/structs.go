package main

type RegisteredMovie struct {
	Id      int
	Imdb_id string
	Title   string
	Year    string
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

type UserRole struct {
	Id     int
	UserId string
	RoleId string
	Key    string
}

type UsersMovie struct {
	Id              int
	User            User
	RegisteredMovie RequestedMovie
	Movie_width     string
	Movie_height    string
	Video_codac     string
	Audio_codac     string
	Container       string
	Frame_rate      string
	Aspect_ratio    string
}

func (userMovie *UsersMovie) OK() error {
	if len(userMovie.Movie_width) == 0 {
		return ErrRequired("Movie Width")
	}
	return nil
}
