// build !release
package button

import "embed"

//go:embed examples/*.html
var examplesfs embed.FS

func (button *Button) Examples() embed.FS {
	return examplesfs
}
