package operation

type Operation struct {
	Name    string
	Execute func() error
}
