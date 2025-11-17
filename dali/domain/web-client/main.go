package webclientdomain

import (
	"strconv"
	"strings"
)

type Domain struct {
	// Define domain fields here
	Path string
}

// this should be wrapping in a specific identifier
func NewDomain() *Domain {
	domain := &Domain{}
	// Initialize Domain fields here
	return domain
}
func (d *Domain) MakePath(keys ...string) string {
	return d.Path + "/" + strings.Join(keys, "/")
}

func (d *Domain) WrapHTML(data []byte, path string, id int64) []byte {
	head := []byte("<x-domain path=\"" + d.Path + "/" + path + "/" + strconv.FormatInt(id, 10) + "\">")
	tail := []byte("</x-domain>")
	result := make([]byte, 0, len(head)+len(data)+len(tail))
	result = append(result, head...)
	result = append(result, data...)
	result = append(result, tail...)
	return result
}
