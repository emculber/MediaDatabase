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

func (task *Task) getTaskWithIdFromDatabase() error {
	log.WithFields(log.Fields{
		"ID": task.Id,
	}).Info("Getting Task With Paramaters")
	var due int64
	err := db.QueryRow(`SELECT id, name, completed, EXTRACT(EPOCH FROM date_trunc('second', due))::INTEGER FROM task_list WHERE id=$1`, task.Id).Scan(&task.Id, &task.Name, &task.Completed, &due)
	task.Due = time.Unix(due, 0) //TODO: DO I NEED THIS
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"ID":   task.Id,
		"Task": task,
	}).Info("Task Found With Paramaters")
	return nil
}

func getTasksFromDatabase() []Task {
	task.Info("Getting Tasks")
	statement := fmt.Sprintf("SELECT id, name, completed, EXTRACT(EPOCH FROM date_trunc('second', due))::INTEGER FROM task_list")
	//TODO: Error Checking
	tasks, _, _ := postgresql_access.QueryDatabase(db, statement)
	task_list := []Task{}

	for _, task := range tasks {
		single_task := Task{}
		single_task.Id, _ = strconv.Atoi(task[0].(string))
		single_task.Name = task[1].(string)
		single_task.Completed, _ = strconv.ParseBool(task[2].(string))
		due_unix, _ := strconv.Atoi(task[3].(string))
		single_task.Due = time.Unix(int64(due_unix), 0)
		task_list = append(task_list, single_task)
	}
	return task_list
}
