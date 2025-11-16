package account

type Group struct {
	id          int64
	name        string
	globalId    string
	Users       []int64
	Groups      []int64
	GetAllUsers func() []int64
	HasUser     func(userId int64) bool
	HasSubGroup func(groupId int64) bool
	IsMemberOf  func() []int64
}

func NewGroup(id int64, name string) *Group {
	return &Group{
		id:       id,
		name:     name,
		globalId: "89",
	}
}

func (u *Group) Id() int64 {
	return u.id
}

func (u *Group) Name() string {
	return u.name
}
