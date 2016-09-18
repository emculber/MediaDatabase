package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type MuxRouter struct {
	Router *mux.Router
}

func init() {
	log.Info("API Init Process Start")
	InitLogger()
	InitDatabase()
	InitExternalSources()
	//InitSecurity()
	log.Info("API Init Process End")
}

func main() {
	log.Info("API Main Mux Router Process Start")
	muxRouter := MuxRouter{}
	muxRouter.Router = mux.NewRouter().StrictSlash(true)
	muxRouter.GenericRouter()
	muxRouter.AdminRouter()
	muxRouter.UserRouter()
	muxRouter.StockRouter()
	muxRouter.TaskRouter()
	log.Info("API Main Mux Router Process End")
	log.Info("API Listening To Port 8080")
	http.ListenAndServe(":8080", muxRouter.Router)
	log.Info("Cloasing DB Connection")
	db.Close()
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
