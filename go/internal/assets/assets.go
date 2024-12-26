// This slurps the data.json file from the assets package directory into the binary at
// build time. This is to support the embedded backend type	for the dota db (data is baked into the
// binary at build time.
//
// If we want to bake other files into the db, we can add them to the assets package directory (like)
// at build time and they'll get slurped up into Assets. You then pass the desired filename to the
// NewEmbeddedDataClient function to load the data from the embedded file.
//
// The data is then served from memory.

package assets

import (
	"embed"
)

//go:embed *
var Assets embed.FS
