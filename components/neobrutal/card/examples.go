// build !release
package card

import "embed"

//go:embed examples/*.html
var examplesfs embed.FS

func (card *Card) Examples() embed.FS {
	return examplesfs
}
