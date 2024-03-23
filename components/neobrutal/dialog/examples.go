// build !release
package dialog

import "embed"

//go:embed examples/*.html
var examplesfs embed.FS

func (dialog *Dialog) Examples() embed.FS {
	return examplesfs
}
