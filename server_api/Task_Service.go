package main

import "fmt"

func (taskTree *TaskTree) taskPlacement(task Task) {
	if task.ParentId == taskTree.Task.Id {
		childrenTaskTree := TaskTree{}
		childrenTaskTree.Task = task
		taskTree.Children = append(taskTree.Children, childrenTaskTree)
		return
	}
	for i, children := range taskTree.Children {
		children.taskPlacement(task)
		taskTree.Children[i] = children
	}
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
