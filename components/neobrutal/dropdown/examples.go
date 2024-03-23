// build !release
package dropdown

import "embed"

//go:embed examples/*.html
var examplesfs embed.FS

func (dropdown *Dropdown) Examples() embed.FS {
	return examplesfs
}
