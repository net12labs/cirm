package job

type JobUnit struct {
	Id     string
	Name   string
	Status string
}

type JobGroup struct {
	*JobUnit
}

type JobSubGroup struct {
	*JobUnit
}

type Job struct {
	*JobUnit
}

type SubJob struct {
	*JobUnit
	ParentId string
}

func NewJob() *JobUnit {
	return &JobUnit{}
}
func NewJobGroup() *JobGroup {
	return &JobGroup{
		JobUnit: &JobUnit{},
	}
}
func NewSubJob() *SubJob {
	return &SubJob{
		JobUnit: &JobUnit{},
	}
}
func NewJobSubGroup() *JobSubGroup {
	return &JobSubGroup{
		JobUnit: &JobUnit{},
	}
}
