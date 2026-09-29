// Package assets embeds extracted production payloads, not original references.
package assets

import "embed"

// Files holds original data banks read by native Go effects.
//
//go:embed raw/*.bin raw/*.mod
var Files embed.FS
