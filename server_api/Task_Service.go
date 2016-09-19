package main

import "fmt"

func (taskTree *TaskTree) taskPlacement(task Task) {
	fmt.Println("Checking if child of ->", taskTree.Task)
	if task.ParentId == taskTree.Task.Id {
		fmt.Println("Task Parent Found")
		childrenTaskTree := TaskTree{}
		childrenTaskTree.Task = task
		taskTree.Children = append(taskTree.Children, childrenTaskTree)
		fmt.Println("Task Tree ->", taskTree)
		fmt.Println("Task Tree Child ->", taskTree.Children)
		return
	}
	fmt.Println("Looping though taskTree.Children ->", taskTree.Children)
	for i, children := range taskTree.Children {
		children.taskPlacement(task)
		taskTree.Children[i] = children
	}
	fmt.Println("Final Task Tree ->", taskTree)
}

func (taskTree *TaskTree) PrintTaskTree(depth int) {
	var indent string
	for i := 0; i < depth; i++ {
		indent += "\t"
	}
	fmt.Println(indent, "ID ->", taskTree.Task.Id)
	fmt.Println(indent, "Name ->", taskTree.Task.Name)
	fmt.Println(indent, "Complete ->", taskTree.Task.Completed)
	fmt.Println(indent, "Due ->", taskTree.Task.Due)

	for _, taskTreeChild := range taskTree.Children {
		taskTreeChild.PrintTaskTree(depth + 1)
	}
}
