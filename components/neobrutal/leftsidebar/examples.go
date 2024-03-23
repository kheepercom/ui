// build !release
package leftsidebar

import "embed"

//go:embed examples/*.html
var examplesfs embed.FS

func (*LeftSidebar) Examples() embed.FS {
	return examplesfs
}
