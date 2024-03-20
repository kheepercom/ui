package ui

// Attributes represents the key-value pairs in an HTML element.
type Attributes map[string][]string

// Get returns the last value associated with the key.
func (a Attributes) Get(key string) string {
	if len(a[key]) == 0 {
		return ""
	}
	return a[key][len(a[key])-1]
}

// GetOr returns the attribute value if set, otherwise returns the default from the 2nd argument.
func (a Attributes) GetOr(key, def string) string {
	v, ok := a.GetHas(key)
	if ok {
		return v
	}
	return def
}

// GetHas returns the last value associated with the given key.
func (a Attributes) GetHas(key string) (string, bool) {
	l := len(a[key])
	if l == 0 {
		return "", false
	}
	return a[key][l-1], true
}

// Has is used with boolean attributes.
func (a Attributes) Has(key string) bool {
	_, ok := a[key]

	return ok
}
