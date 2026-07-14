package bytescty

import _ "embed"

//go:embed externs.cty
var externsCty []byte

// ExternsFilename is the name reported for the embedded declarations in
// diagnostics.
const ExternsFilename = "bytes-cty-type/externs.cty"

// Externs returns the functy `//functy:extern` declarations for the functions
// GetBytesFunctions provides: their real signatures, which their cty metadata cannot
// express.
//
// All three take an argument that is a union — a string, or a bytes value — and cty
// has no union type, so its metadata can only say "dynamic". Two of them also take an
// optional trailing content type, and cty can only make an argument optional by making
// it variadic, which erases its name, its type, and its arity; for base64decode, whose
// return type is a string or a bytes value depending on whether that argument is
// present, it erases the signature entirely. These declarations say what each really
// accepts, so that help(), generated documentation, and editor tooling can show it.
//
// The bytes are opaque to this package: it does not import functy, and nothing here
// parses them. A functy host registers them:
//
//	parser.RegisterExterns(bytescty.Externs(), bytescty.ExternsFilename)
func Externs() []byte { return externsCty }
