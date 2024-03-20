// build !release
package checkbox

import "embed"

//go:embed examples/*.html
var examplesfs embed.FS

func (checkbox *Checkbox) Examples() embed.FS {
	return examplesfs
}
