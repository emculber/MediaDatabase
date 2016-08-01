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

	InitLogger()
	InitDatabase()
	InitExternalSources()
	//InitSecurity()

}

func main() {
	router := NewRouter()
	http.ListenAndServe(":8080", router)
}

func respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	var buffer bytes.Buffer
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
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}

	return nil
}
