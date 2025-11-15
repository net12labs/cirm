package service

type RunMode struct {
	items map[string]string
}

func (r *RunMode) SetKey(key string) {
	if r.items == nil {
		r.items = make(map[string]string)
	}
	r.items[key] = ""
}

func (r *RunMode) SetKeys(keys ...string) {
	if r.items == nil {
		r.items = make(map[string]string)
	}
	for _, key := range keys {
		r.items[key] = ""
	}
}

func (r *RunMode) SetKV(key, value string) {
	if r.items == nil {
		r.items = make(map[string]string)
	}
	r.items[key] = value
}

func (r *RunMode) HasKey(key string) bool {
	_, exists := r.items[key]
	return exists
}
func (r *RunMode) GetValue(key string) string {
	return r.items[key]
}
func (r *RunMode) HasValue(key, value string) bool {
	if v, exists := r.items[key]; exists {
		return v == value
	}
	return false
}
