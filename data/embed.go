package data

import "embed"

// Files holds default game content shipped in the binary.
//
//go:embed buildings.json events.json
var Files embed.FS
