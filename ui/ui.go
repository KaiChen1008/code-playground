package ui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var static embed.FS

// Static is the root of the UI assets (the contents of the dist/ directory)
var Static, _ = fs.Sub(static, "dist")
