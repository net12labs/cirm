package unit

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Listeners struct {
	items []func(any)
}

func (l *Listeners) AddListener(fn func(any)) *Listeners {
	l.items = append(l.items, fn)
	return l
}
func (l *Listeners) RemoveListener(fn func(any)) {
	for i, v := range l.items {
		if &v == &fn {
			l.items = append(l.items[:i], l.items[i+1:]...)
			break
		}
	}
}

func dispatchListeners(listeners *Listeners, arg any) {
	for _, listener := range listeners.items {
		listener(arg)
	}
}

type RuntimeUnit struct {
	OnPanic *Listeners
	OnExit  *Listeners
	// RuntimeUnit fields here
}

func NewRuntimeUnit() *RuntimeUnit {
	return &RuntimeUnit{
		OnPanic: &Listeners{},
		OnExit:  &Listeners{},
	}
}

func (r *RuntimeUnit) Panic(err string) {
	dispatchListeners(r.OnPanic, err)
	panic(err)

}

func (r *RuntimeUnit) Exit(code int) {
	dispatchListeners(r.OnExit, code)
	os.Exit(code)
}

func (r *RuntimeUnit) ExitErr(code int, err error) {
	fmt.Println("Exiting with error:", err)
	dispatchListeners(r.OnExit, code)
	os.Exit(code)
}

func (r *RuntimeUnit) WaitForSIGTERM() {
	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	r.Exit(0)
}
