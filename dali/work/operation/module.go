package operation

type OperationUnit struct {
	Name    string
	Execute func() error
}

type Operation struct {
	*OperationUnit
}

type SubOperation struct {
	*OperationUnit
	ParentId string
}
type OperationGroup struct {
	*OperationUnit
}

type OperationSubGroup struct {
	*OperationUnit
	ParentId string
}

func NewOperation() *Operation {
	return &Operation{
		OperationUnit: &OperationUnit{},
	}
}

func NewOperationGroup() *OperationGroup {
	return &OperationGroup{
		OperationUnit: &OperationUnit{},
	}
}

func NewSubOperation() *SubOperation {
	return &SubOperation{
		OperationUnit: &OperationUnit{},
	}
}

func NewOperationSubGroup() *OperationSubGroup {
	return &OperationSubGroup{
		OperationUnit: &OperationUnit{},
	}
}
