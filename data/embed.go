package data

import "embed"

// Files holds default game content shipped in the binary.
//
//go:embed buildings.json events.json scenarios.json quiet_beats.json
var Files embed.FS
