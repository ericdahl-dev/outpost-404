package data

import "embed"

// Files holds default game content shipped in the binary.
//
//go:embed buildings.json events.json scenarios.json
var Files embed.FS
