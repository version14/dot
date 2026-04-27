package generator

import "github.com/version14/dot/pkg/dotapi"

// Manifest, Command, Validator, Check, CheckType are defined in pkg/dotapi so
// generator authors can declare them without depending on internal packages.
// The internal generator runtime imports the canonical types via these aliases
// to avoid circular references.
type (
	Manifest  = dotapi.Manifest
	Command   = dotapi.Command
	Validator = dotapi.Validator
	Check     = dotapi.Check
	CheckType = dotapi.CheckType
)

const (
	CheckFileExists    = dotapi.CheckFileExists
	CheckJSONKeyExists = dotapi.CheckJSONKeyExists
)
