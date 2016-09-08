package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func newTask(w http.ResponseWriter, r *http.Request) {
	task := Task{}
	fmt.Println(r.Body)
	//buf := new(bytes.Buffer)
	//buf.ReadFrom(r.Body)
	//s := buf.String() // Does a complete copy of the bytes in the buffer.
	//fmt.Println(s)
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), 400)
		fmt.Println(err)
		return
	}

	fmt.Println(task.RegisterNewTask())
	fmt.Println("Task ->", task)

	w.Write([]byte("Task Was Created"))
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(getTasksFromDatabase()); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error Encoding Wallet")
		return
	}
}
