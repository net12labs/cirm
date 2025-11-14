package xmodule

func NewEtcStore() *EtcStore {
	store := EtcStore{
		Data:        map[string]string{},
		subscribers: []*EtcValue{},
	}
	return &store
}

func NewEtcStoreCb(callback func(*EtcStore)) *EtcStore {
	store := NewEtcStore()
	callback(store)
	return store
}
