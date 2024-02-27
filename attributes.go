package ui

// Attributes represents the key-value pairs in an HTML element.
type Attributes map[string][]string

// Get returns the last value associated with the given key.
func (a Attributes) Get(key string) (string, bool) {
	l := len(a[key])
	if l == 0 {
		return "", false
	}
	return a[key][l-1], true
}
