package work

type Work struct {
	// Work fields here
}

type WorkUnit struct {
	*Work
}

type WorkGroup struct {
	*Work
}

type WorkSubUnit struct {
	*Work
	ParentId string
}

type WorkSubGroup struct {
	*Work
	ParentId string
}

func NewWorkUnit() *WorkUnit {
	return &WorkUnit{
		Work: &Work{},
	}
}

func NewWorkGroup() *WorkGroup {
	return &WorkGroup{
		Work: &Work{},
	}
}

func NewWorkSubUnit(parentId string) *WorkSubUnit {
	return &WorkSubUnit{
		Work:     &Work{},
		ParentId: parentId,
	}
}

func NewWorkSubGroup(parentId string) *WorkSubGroup {
	return &WorkSubGroup{
		Work:     &Work{},
		ParentId: parentId,
	}
}
