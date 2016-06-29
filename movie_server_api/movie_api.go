package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

func init() {

	//InitLogger()
	//InitDatabase()
	//InitExternalSources()

}

type Adapter func(http.Handler) http.Handler

func ApiCallLog() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//api_logger_fields := ApiLoggerFields{}
			//api_logger_fields.ip_address = r.RemoteAddr
			//api_logger_fields.method_type = r.Method
			log.WithFields(log.Fields{
				"Logger Fields": r.RemoteAddr,
			}).Info("API Call")
			h.ServeHTTP(w, r)
		})
	}
}

//func ValidateUserKey(user User) Adapter {
func ValidateUserKey() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.WithFields(log.Fields{
				"User": "test",
			}).Info("Validating User")
			h.ServeHTTP(w, r)
		})
	}
}

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

/*
func addMovieToUserMovies(w http.ResponseWriter, r *http.Request) {
	api_logger_fields := ApiLoggerFields{}
	api_logger_fields.ip_address = r.RemoteAddr
	api_logger_fields.method_type = r.Method

	if r.Method != "POST" {
		log.WithFields(log.Fields{
			"Logger Fields": api_logger_fields,
			"Error":         http.StatusMethodNotAllowed,
		}).Error("Invalid Request!")
		http.Error(w, "Invalid Request -> Incorrect Method Call", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	r.ParseForm()

	user := User{}
	user.User_key = r.PostFormValue("user_key")

	omdbapi := OmdbapiData{}
	omdbapi.Imdb_id = strings.ToLower(r.PostFormValue("imdb_id"))

	users_movie := UsersMovie{}
	users_movie.User = user
	users_movie.Omdbapi = omdbapi
	users_movie.Movie_width = r.PostFormValue("movie_width")
	users_movie.Movie_height = r.PostFormValue("movie_height")
	users_movie.Video_codac = r.PostFormValue("video_codac")
	users_movie.Audio_codac = r.PostFormValue("audio_codac")
	users_movie.Container = r.PostFormValue("container")
	users_movie.Frame_rate = r.PostFormValue("frame_rate")
	users_movie.Aspect_ratio = r.PostFormValue("aspect_ratio")

	if users_movie.User.User_key == "" || users_movie.Omdbapi.Imdb_id == "" || users_movie.Movie_width == "" ||
		users_movie.Movie_height == "" || users_movie.Audio_codac == "" || users_movie.Video_codac == "" ||
		users_movie.Container == "" || users_movie.Frame_rate == "" || users_movie.Aspect_ratio == "" {

		log.WithFields(log.Fields{
			"Logger Fields": api_logger_fields,
			"Error":         "addMovieToUserList -> Empty Value Detected",
			"User Movie":    users_movie,
		}).Error("Empty Content")
		http.Error(w, "Invalid Request -> Empty Content Was Detected", http.StatusBadRequest)
		return
	}

	var isValidated bool
	var err error
	isValidated, users_movie.User, err = validateUserId(users_movie.User)
	if err != nil || !isValidated {
		log.WithFields(log.Fields{
			"Logger Fields":    api_logger_fields,
			"Error":            http.StatusBadRequest,
			"Validation Error": err,
			"User":             users_movie.User,
			"Data":             users_movie,
			"Registered User":  isValidated,
		}).Error("Invalid User")
		http.Error(w, "Invalid Request -> Invalid User", http.StatusBadRequest)
		return
	}

	isValidated, users_movie.Omdbapi, err = validateImdbId(users_movie.Omdbapi)
	if err != nil || !isValidated {
		users_movie.Omdbapi, err = Omdbapi(users_movie.Omdbapi.Imdb_id)
		if err != nil {
			log.WithFields(log.Fields{
				"Logger Fields": api_logger_fields,
				"Error":         err,
				"Movie":         users_movie.Omdbapi,
				"Data":          users_movie,
			}).Error("Error Finding Movie")
		}
		users_movie.Omdbapi, err = InsertNewMovie(users_movie.Omdbapi)
		if err != nil {
			log.WithFields(log.Fields{
				"Logger Fields": api_logger_fields,
				"Error":         err,
				"Movie":         users_movie.Omdbapi,
				"Data":          users_movie,
			}).Error("Error Inserting Movie")
		}
		log.WithFields(log.Fields{
			"Validation Error": err,
			"Movie":            users_movie.Omdbapi,
			"Data":             users_movie,
			"Registered Movie": isValidated,
		}).Info("Added movie")
	}

	log.WithFields(log.Fields{
		"Data": users_movie,
	}).Info("Adding Users Movie")

	isValidated, err = validateUsersMovie(users_movie)
	if isValidated {
		log.WithFields(log.Fields{
			"Logger Fields": api_logger_fields,
			"Data":          users_movie,
		}).Info("Movie Already Added")
	} else {
		users_movie, err = InsertNewUsersMovie(users_movie)
		if err != nil {
			log.WithFields(log.Fields{
				"Logger Fields":  api_logger_fields,
				"Error":          http.StatusBadRequest,
				"Database Error": err,
				"Data":           users_movie,
			}).Error("Error adding data to database")
			http.Error(w, "Invalid Request -> Error Adding Movie", http.StatusBadRequest)
			return
		}
		log.WithFields(log.Fields{
			"Data": users_movie,
		}).Info("Movie Added")
	}
	w.Write([]byte("OK"))
}

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	api_logger_fields := ApiLoggerFields{}
	api_logger_fields.ip_address = r.RemoteAddr
	api_logger_fields.method_type = r.Method

	if r.Method != "POST" {
		log.WithFields(log.Fields{
			"Logger Fields": api_logger_fields,
			"Error":         http.StatusMethodNotAllowed,
		}).Error("Invalid Request!")
		http.Error(w, "Invalid Request!", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()

	user := User{}
	user.User_key = r.PostFormValue("user_key")

	var isValidated bool
	var err error
	isValidated, user, err = validateUserId(user)
	if err != nil || !isValidated {
		log.WithFields(log.Fields{
			"Logger Fields":    api_logger_fields,
			"Error":            http.StatusBadRequest,
			"Validation Error": err,
			"User":             user,
			"Registered User":  isValidated,
		}).Error("Invalid User")
		http.Error(w, "Invalid Request -> Invalid User", http.StatusBadRequest)
		return
	}

	movies := ReadUserMovies(user)

	if err = json.NewEncoder(w).Encode(movies); err != nil {
		log.WithFields(log.Fields{
			"Logger Fields": api_logger_fields,
			"Error":         err,
		}).Error("Invalid Request!")
		panic(err)
	}
}

*/
func main() {
	//http.HandleFunc("/api/getallmovies", getAllMovies)
	//http.HandleFunc("/api/addmovie", addMovieToUserMovies)
	http.Handle("/api/test", Adapt(test(), ValidateUserKey(), ApiCallLog()))
	http.ListenAndServe(":8080", nil)
}

func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	var buffer bytes.Buffer
	var err error
	if err := json.NewEncoder(&buffer).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	if _, err := io.Copy(w, &buffer); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Respond Error ->")
	}
}

func decode(r *http.Request, data interface{}) error {
	var err error
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}

	//if err = data.OK(); err != nil {
	//		return err
	//	}

	return nil
}
