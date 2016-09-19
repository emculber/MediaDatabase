package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func newTask(w http.ResponseWriter, r *http.Request) {
	task := Task{}
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(r.Body)
	//s := buf.String() // Does a complete copy of the bytes in the buffer.
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Decoding Json")
		return
	}

	task.RegisterNewTask()
	log.WithFields(log.Fields{
		"Task": task,
	}).Info("Regerstering New Task")

	w.Write([]byte("Task Was Created"))
}

func getTask(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting Task")
	taskTree := TaskTree{}
	err := json.NewDecoder(r.Body).Decode(&taskTree.Task)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Decoding Json")
		return
	}

	taskTree.getTaskWithIdFromDatabase()

	fmt.Println(taskTree)

	if err := json.NewEncoder(w).Encode(taskTree); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Encoding Wallet")
		return
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(getTasksFromDatabase()); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Encoding Wallet")
		return
	}
}
