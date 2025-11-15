package task

type TaskUnit struct {
	Id     string
	Name   string
	Status string
}

type TaskGroup struct {
	*TaskUnit
}
type TaskSubGroup struct {
	*TaskUnit
	ParentId string
}

type Task struct {
	*TaskUnit
}

type SubTask struct {
	*TaskUnit
	ParentId string
}

func NewTask() *Task {
	return &Task{
		TaskUnit: &TaskUnit{},
	}
}

func NewTaskGroup() *TaskGroup {
	return &TaskGroup{
		TaskUnit: &TaskUnit{},
	}
}

func NewSubTask() *SubTask {
	return &SubTask{
		TaskUnit: &TaskUnit{},
	}
}

func NewTaskSubGroup() *TaskSubGroup {
	return &TaskSubGroup{
		TaskUnit: &TaskUnit{},
	}
}
