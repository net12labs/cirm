package xmodule

import "fmt"

type EtcStore struct {
	Data        map[string]string
	subscribers []*EtcValue
	locks       map[string]bool
}

func (e *EtcStore) Subscribe(val *EtcValue) {
	for _, sub := range e.subscribers {
		if sub == val {
			return // Already subscribed
		}
	}
	e.subscribers = append(e.subscribers, val)
}
func (e *EtcStore) Unsubscribe(val *EtcValue) {
	for i, sub := range e.subscribers {
		if sub == val {
			e.subscribers = append(e.subscribers[:i], e.subscribers[i+1:]...)
			break
		}
	}
}

func (e *EtcStore) Lock(key string) {
	e.locks[key] = true
	e.dispatchLocked(key)
}
func (e *EtcStore) Unlock(key string) {
	delete(e.locks, key)
	e.dispatchUnlocked(key)
}

func (e *EtcStore) NewValue(key string) *EtcValue {
	return &EtcValue{
		Key:      key,
		etcStore: e,
	}
}

func (e *EtcStore) ErrorLocked(key string) error {
	if e.locks[key] {
		return fmt.Errorf("key %q is locked", key)
	}
	return nil
}

func (e *EtcStore) Get(key string) *EtcValue {
	value, exists := e.Data[key]
	val := e.NewValue(key)
	val.Value = value
	val.exists = exists
	return val
}

func (e *EtcStore) SetKV(key, value string) error {
	if err := e.ErrorLocked(key); err != nil {
		return err
	}
	v, exists := e.Data[key]
	e.Data[key] = value
	if !exists {
		e.dispatchAppeared(key, value)
	}

	if v != value {
		e.dispatchChanged(key, value)
	}
	return nil
}

func (e *EtcStore) Set(val *EtcValue) error {
	if err := e.ErrorLocked(val.Key); err != nil {
		return err
	}
	v, exists := e.Data[val.Key]
	e.Data[val.Key] = val.Value
	if !exists {
		e.dispatchAppeared(val.Key, val.Value)
	}

	if v != val.Value {
		e.dispatchChanged(val.Key, val.Value)
	}
	return nil
}

func (e *EtcStore) Delete(key string) error {
	if err := e.ErrorLocked(key); err != nil {
		return err
	}
	_, exists := e.Data[key]
	if !exists {
		return fmt.Errorf("does not exist| key %q does not exist", key)
	}
	delete(e.Data, key)
	e.dispatchLeft(key)
	return nil
}

func (e *EtcStore) dispatchChanged(key string, val string) {
	// Notify subscribers about the change
	for _, sub := range e.subscribers {
		if sub.Key == key && sub.OnChanged != nil {
			sub.OnChanged(val)
		}
	}
}

func (e *EtcStore) dispatchLeft(key string) {
	// Notify subscribers about the deletion
	for _, sub := range e.subscribers {
		if sub.Key == key && sub.OnLeft != nil {
			sub.OnLeft()
		}
	}
}
func (e *EtcStore) dispatchAppeared(key string, val string) {
	// Notify subscribers about the appearance
	for _, sub := range e.subscribers {
		if sub.Key == key && sub.OnAppeared != nil {
			sub.OnAppeared(val)
		}
	}
}
func (e *EtcStore) dispatchUnlocked(key string) {
	// Notify subscribers about the unlocking
	for _, sub := range e.subscribers {
		if sub.Key == key && sub.OnUnlocked != nil {
			sub.OnUnlocked()
		}
	}
}
func (e *EtcStore) dispatchLocked(key string) {
	// Notify subscribers about the locking
	for _, sub := range e.subscribers {
		if sub.Key == key && sub.OnLocked != nil {
			sub.OnLocked()
		}
	}
}
