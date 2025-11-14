package webserver

import (
	"fmt"
	"strings"
)

type URLPath struct {
	Path string
}

func (up *URLPath) SegmentIsEqual(index, key string) bool {
	return up.Path == fmt.Sprintf("/%s/%s", index, key)
}
func (up *URLPath) IsEqual(path string) bool {
	return up.Path == path
}
func (up *URLPath) GetSegment(index int) string {
	segments := strings.Split(up.Path, "/")
	if index < 0 || index >= len(segments) {
		return ""
	}
	return segments[index]
}
func (up *URLPath) StartsWith(prefix string) bool {
	return strings.HasPrefix(up.Path, prefix)
}
