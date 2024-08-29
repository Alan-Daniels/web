package internal

import "embed"

//go:embed commit.txt
var Commit string

var unused embed.FS
