// filepath: /home/lxk/Desktop/cirm/bins/vfsql/events.go
package vfsql

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Subscribe creates a new event subscription
func (vfs *VFS) Subscribe(filter *EventFilter) (*EventSubscription, error) {
	if filter == nil {
		filter = &EventFilter{}
	}

	// Set default buffer size
	if filter.BufferSize == 0 {
		filter.BufferSize = 100
	}

	sub := &EventSubscription{
		ID:     fmt.Sprintf("sub-%d-%d", vfs.id, currentTimestamp()),
		Events: make(chan Event, filter.BufferSize),
		Errors: make(chan error, 10),
		filter: filter,
		cancel: make(chan struct{}),
	}

	vfs.subMutex.Lock()
	vfs.subscribers[sub.ID] = sub
	vfs.subMutex.Unlock()

	return sub, nil
}

// Unsubscribe removes an event subscription
func (vfs *VFS) Unsubscribe(sub *EventSubscription) error {
	vfs.subMutex.Lock()
	defer vfs.subMutex.Unlock()

	if _, ok := vfs.subscribers[sub.ID]; !ok {
		return nil // Already unsubscribed
	}

	delete(vfs.subscribers, sub.ID)
	close(sub.cancel)
	close(sub.Events)
	close(sub.Errors)

	return nil
}

// emitEvent sends an event to all matching subscribers
func (vfs *VFS) emitEvent(event Event) {
	vfs.subMutex.RLock()
	subscribers := make([]*EventSubscription, 0, len(vfs.subscribers))
	for _, sub := range vfs.subscribers {
		subscribers = append(subscribers, sub)
	}
	vfs.subMutex.RUnlock()

	// Dispatch to matching subscribers in goroutines
	for _, sub := range subscribers {
		go vfs.dispatchEvent(sub, event)
	}
}

// dispatchEvent sends an event to a subscriber if it matches the filter
func (vfs *VFS) dispatchEvent(sub *EventSubscription, event Event) {
	// Check if subscription is cancelled
	select {
	case <-sub.cancel:
		return
	default:
	}

	// Check if event matches filter
	if !vfs.matchesFilter(sub.filter, event) {
		return
	}

	// Try to send event (non-blocking)
	select {
	case sub.Events <- event:
		// Event sent successfully
	case <-sub.cancel:
		// Subscription cancelled
	default:
		// Channel full, send error
		select {
		case sub.Errors <- fmt.Errorf("event buffer full, dropping event"):
		default:
			// Error channel also full, nothing we can do
		}
	}
}

// matchesFilter checks if an event matches a subscription filter
func (vfs *VFS) matchesFilter(filter *EventFilter, event Event) bool {
	// Check event type filter
	if len(filter.EventTypes) > 0 {
		matched := false
		for _, et := range filter.EventTypes {
			if et == event.Type {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Check file type filter
	if len(filter.FileTypes) > 0 {
		matched := false
		for _, ft := range filter.FileTypes {
			if ft == event.FileType {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Check path filter
	if len(filter.Paths) > 0 {
		matched := false
		for _, filterPath := range filter.Paths {
			if vfs.pathMatches(filterPath, event.Path, filter.Recursive) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Check name pattern
	if filter.NamePattern != "" {
		name := filepath.Base(event.Path)
		matched, err := filepath.Match(filter.NamePattern, name)
		if err != nil || !matched {
			return false
		}
	}

	// Check tag filter (only for events that have tags)
	if len(filter.Tags) > 0 && len(event.Tags) > 0 {
		if filter.TagMatchAll {
			// All filter tags must be present
			for _, filterTag := range filter.Tags {
				found := false
				for _, eventTag := range event.Tags {
					if eventTag == filterTag {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		} else {
			// At least one filter tag must be present
			matched := false
			for _, filterTag := range filter.Tags {
				for _, eventTag := range event.Tags {
					if eventTag == filterTag {
						matched = true
						break
					}
				}
				if matched {
					break
				}
			}
			if !matched {
				return false
			}
		}
	}

	return true
}

// pathMatches checks if an event path matches a filter path
func (vfs *VFS) pathMatches(filterPath, eventPath string, recursive bool) bool {
	filterPath = normalizePath(filterPath)
	eventPath = normalizePath(eventPath)

	if filterPath == eventPath {
		return true
	}

	if recursive {
		// Check if eventPath is under filterPath
		return strings.HasPrefix(eventPath, filterPath+"/") || filterPath == "/"
	}

	// Non-recursive: check if parent directory matches
	eventDir, _ := splitPath(eventPath)
	return eventDir == filterPath
}
