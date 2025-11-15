package etcstore

import (
	"fmt"
	"strconv"
)

type EtcValue struct {
	etcStore   *EtcStore
	Key        string
	Value      string
	exists     bool
	OnChanged  func(newValue string)
	OnLeft     func()
	OnAppeared func(val string)
	OnUnlocked func()
	OnLocked   func()
}

func (e *EtcValue) IsLocked() bool {
	_, locked := e.etcStore.locks[e.Key]
	return locked
}

func (e *EtcValue) Exists() bool {
	return e.exists
}
func (e *EtcValue) ErrNotExists() error {
	if !e.exists {
		return fmt.Errorf("key %q does not exist", e.Key)
	}
	return nil
}

func (e *EtcValue) Watch() {
	e.etcStore.Subscribe(e)
}
func (e *EtcValue) Unwatch() {
	e.etcStore.Unsubscribe(e)
}

func (e *EtcValue) Refresh() {
	e.etcStore.Data[e.Key] = e.Value
}
func (e *EtcValue) HasChanged() bool {
	return e.etcStore.Data[e.Key] != e.Value
}
func (e *EtcValue) RefreshIfChanged() {
	if e.HasChanged() {
		e.Refresh()
	}
}

// casts //

func (e *EtcValue) String() string {
	return e.Value
}
func (e *EtcValue) IsBool() bool {
	_, err := strconv.ParseBool(e.Value)
	return err == nil
}

func (e *EtcValue) Bool() bool {
	value, _ := strconv.ParseBool(e.Value)
	return value
}

func (e *EtcValue) IsInt() bool {
	_, err := strconv.Atoi(e.Value)
	return err == nil
}
func (e *EtcValue) Int() int {
	value, _ := strconv.Atoi(e.Value)
	return value
}
