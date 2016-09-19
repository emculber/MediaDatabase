package main

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/emculber/database_access/postgresql"
)

var taskDatabaseSchema = []string{
	"CREATE TABLE task_list(id SERIAL PRIMARY KEY, name VARCHAR(60), completed BOOLEAN, due TIMESTAMP, parent INTEGER REFERENCES task_list(id))",
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
	if task.ParentId != 0 {
		err := db.QueryRow(`insert into task_list (name, completed, due, parent) values($1, $2, $3, $4) returning id`, task.Name, task.Completed, task.Due, task.ParentId).Scan(&task.Id)
		fmt.Println(err)
		if err != nil {
			return err
		}
	} else {
		err := db.QueryRow(`insert into task_list (name, completed, due) values($1, $2, $3) returning id`, task.Name, task.Completed, task.Due).Scan(&task.Id)
		fmt.Println(err)
		if err != nil {
			return err
		}
	}
	return nil
}

func (taskTree *TaskTree) getTaskWithIdFromDatabase() error {
	log.WithFields(log.Fields{
		"ID": taskTree.Task.Id,
	}).Info("Getting Task With Paramaters")
	statement := fmt.Sprintf(`WITH RECURSIVE parent_task(id, name, completed, due, parent) AS (
														  SELECT
															  id,
															  name,
															  completed,
															  EXTRACT(EPOCH FROM date_trunc('second', due))::INTEGER,
															  parent
														  FROM task_list
															WHERE id = %d
															UNION ALL SELECT
																					ct.id,
																					ct.name,
																					ct.completed,
																					EXTRACT(EPOCH FROM date_trunc('second', ct.due))::INTEGER,
																					ct.parent
																				FROM parent_task pt, task_list ct
																				WHERE ct.parent = pt.id) SELECT
																																	 id,
																																	 name,
																																	 completed,
																																	 due,
																																	 parent
																																 FROM parent_task`, taskTree.Task.Id)
	data, _, err := postgresql_access.QueryDatabase(db, statement)

	fmt.Println(data)

	for _, task := range data {
		single_task := Task{}
		single_task.Id, _ = strconv.Atoi(task[0].(string))
		single_task.Name = task[1].(string)
		single_task.Completed, _ = strconv.ParseBool(task[2].(string))
		due, _ := strconv.Atoi(task[3].(string))
		single_task.Due = time.Unix(int64(due), 0)
		single_task.ParentId, _ = strconv.Atoi(task[4].(string))
		fmt.Println("Finding place for task ->", single_task)
		if single_task.Id == taskTree.Task.Id {
			fmt.Println("Root Parent Task ->", single_task)
			taskTree.Task = single_task
			continue
		}
		taskTree.taskPlacement(single_task)
	}

	taskTree.PrintTaskTree(0)

	//task.Due = time.Unix(due, 0) //TODO: DO I NEED THIS
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"ID":        taskTree.Task.Id,
		"Task Tree": taskTree,
	}).Info("Task Found With Paramaters")
	return nil
}

func (taskTree *TaskTree) getTasksFromDatabase() {
	log.Info("Getting Tasks")
	statement := fmt.Sprintf("SELECT id, name, completed, EXTRACT(EPOCH FROM date_trunc('second', due))::INTEGER, parent FROM task_list")
	//TODO: Error Checking
	fmt.Println(statement)
	tasks, _, err := postgresql_access.QueryDatabase(db, statement)
	fmt.Println(err)

	for _, task := range tasks {
		single_task := Task{}
		single_task.Id, _ = strconv.Atoi(task[0].(string))
		single_task.Name = task[1].(string)
		single_task.Completed, _ = strconv.ParseBool(task[2].(string))
		due_unix, _ := strconv.Atoi(task[3].(string))
		single_task.Due = time.Unix(int64(due_unix), 0)
		single_task.ParentId, _ = strconv.Atoi(task[4].(string))
		fmt.Println("Finding place for task ->", single_task)
		if single_task.Id == taskTree.Task.Id {
			fmt.Println("Root Parent Task ->", single_task)
			taskTree.Task = single_task
			continue
		}
		taskTree.taskPlacement(single_task)
	}
	taskTree.PrintTaskTree(0)
}
