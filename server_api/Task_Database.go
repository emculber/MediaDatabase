package main

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/emculber/database_access/postgresql"
)

var taskDatabaseSchema = []string{
	"CREATE TABLE task_list(id SERIAL PRIMARY KEY, name VARCHAR(60), completed BOOLEAN, due TIMESTAMP)",
}

var taskDropDatabaseSchema = []string{
	"DROP TABLE task_list",
}

func CreateTaskTables() {
	//TODO: check if table exists
	for _, table := range taskDatabaseSchema {
		log.WithFields(log.Fields{
			"Table": table,
		}).Info("Creating Table")
		err := postgresql_access.CreateDatabaseTable(db, table)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error Creating Table")
		}
	}
}

func DropTaskTables() {
	//TODO: check if table exists
	for _, table := range taskDropDatabaseSchema {
		log.WithFields(log.Fields{
			"Table": table,
		}).Info("Drop Table")
		err := postgresql_access.CreateDatabaseTable(db, table)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error Drop Table")
		}
	}
}

func (task *Task) RegisterNewTask() error {
	err := db.QueryRow(`insert into task_list (name, completed, due) values($1, $2, $3) returning id`, task.Name, task.Completed, task.Due).Scan(&task.Id)
	if err != nil {
		return err
	}
	return nil
}

func getTasksFromDatabase() []Task {
	fmt.Println("Getting Tasks")
	statement := fmt.Sprintf("SELECT id, name, completed, due FROM task_list")
	//TODO: Error Checking
	fmt.Println(statement)
	tasks, _, err := postgresql_access.QueryDatabase(db, statement)
	fmt.Println(err)
	fmt.Println(tasks)
	task_list := []Task{}

	for _, task := range tasks {
		single_task := Task{}
		single_task.Id, _ = strconv.Atoi(task[0].(string))
		single_task.Name = task[1].(string)
		single_task.Completed = task[2].(bool)
		single_task.Due = task[3].(time.Time)
		task_list = append(task_list, single_task)
	}
	fmt.Println("Returning tasks ->", len(task_list))
	return task_list
}
