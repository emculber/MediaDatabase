package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	fmt.Println("API Init Process Start")
	InitLogger()
	InitDatabase()
	InitExternalSources()
	//InitSecurity()
	fmt.Println("API Init Process End")
}

func main() {
	fmt.Println("API Main Mux Router Process Start")
	muxRouter := MuxRouter{}
	muxRouter.Router = mux.NewRouter().StrictSlash(true)
	muxRouter.GenericRouter()
	muxRouter.AdminRouter()
	muxRouter.UserRouter()
	muxRouter.StockRouter()
	muxRouter.TaskRouter()
	fmt.Println("API Main Mux Router Process End")
	fmt.Println("API Listening To Port 8080")
	http.ListenAndServe(":8080", muxRouter.Router)
	fmt.Println("Cloasing DB Connection")
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
