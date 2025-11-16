package project

type ProjectUnit struct {
	Id     string
	Name   string
	Status string
}

type ProjectGroup struct {
	*ProjectUnit
}

type Project struct {
	*ProjectUnit
}

type SubProject struct {
	*ProjectUnit
	ParentId string
}

type ProjectSubGroup struct {
	*ProjectUnit
	ParentId string
}

func NewProject() *Project {
	return &Project{
		ProjectUnit: &ProjectUnit{},
	}
}

func NewProjectGroup() *ProjectGroup {
	return &ProjectGroup{
		ProjectUnit: &ProjectUnit{},
	}
}

func NewSubProject() *SubProject {
	return &SubProject{
		ProjectUnit: &ProjectUnit{},
	}
}

func NewProjectSubGroup() *ProjectSubGroup {
	return &ProjectSubGroup{
		ProjectUnit: &ProjectUnit{},
	}
}
