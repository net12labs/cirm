package xmodule

func NewEtcStore() *EtcStore {
	store := EtcStore{
		Data:        map[string]string{},
		subscribers: []*EtcValue{},
	}
	return &store
}
