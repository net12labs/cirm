package task

type Task struct {
	Id   string
	Name string
}

type AgentTask struct {
	Task    Task
	AgentId string
}

type HumanTask struct {
	Task    Task
	HumanId string
}
