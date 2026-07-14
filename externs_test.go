package bytescty

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

// externDeclRE matches a top-level declaration in externs.cty. The file is parsed here
// with a regex rather than with functy on purpose: this package must not depend on
// functy (its bytes are opaque to it), and the check only needs the name set.
var externDeclRE = regexp.MustCompile(`(?m)^func (\w+)\(`)

// TestExternsCoverEveryFunction is the drift guard. Every function this package
// provides is declared in externs.cty — none of the three has a cty signature that
// tells the truth, since each takes a string-or-bytes union that cty can only call
// "dynamic". Adding a function without declaring it would leave it reflecting as that
// lie; declaring one that no longer exists would document a function nobody can call.
// Either way, this fails.
func TestExternsCoverEveryFunction(t *testing.T) {
	declared := make(map[string]bool)
	for _, m := range externDeclRE.FindAllStringSubmatch(string(Externs()), -1) {
		declared[m[1]] = true
	}

	funcs := GetBytesFunctions()
	for name := range funcs {
		assert.True(t, declared[name],
			"%s() is provided by GetBytesFunctions but has no declaration in externs.cty, "+
				"so it reflects as its cty signature — which cannot say that its argument is "+
				"a string or a bytes value, nor how many arguments it takes", name)
	}
	for name := range declared {
		assert.Contains(t, funcs, name,
			"externs.cty declares %s(), which GetBytesFunctions does not provide", name)
	}
}

// The bytes must declare themselves an extern file: functy's RegisterExterns verifies
// the directive rather than forcing the mode, so that this same file is a valid
// standalone .cty that `functy fmt` and `functy symbols` can open.
func TestExternsCarryTheDirective(t *testing.T) {
	require.True(t, strings.HasPrefix(string(Externs()), "//functy:extern\n"),
		"externs.cty must begin with the //functy:extern directive")
}

// Every function, and every parameter of every function, must carry a cty description.
//
// The cty metadata is the only documentation a non-functy cty host can see, and the
// only thing functy's own doc() reads (doc() does not consult the extern), so a gap
// here reads as "exists but undocumented" even where help() shows a full block. An
// extern says what a signature *is*; it does not excuse the metadata from saying what
// the function does.
func TestEverythingIsDescribed(t *testing.T) {
	for name, fn := range GetBytesFunctions() {
		assert.NotEmpty(t, fn.Description(), "%s() has no cty Description", name)

		for _, p := range fn.Params() {
			assert.NotEmpty(t, p.Description, "%s() parameter %q has no Description", name, p.Name)
		}
		if vp := fn.VarParam(); vp != nil {
			assert.NotEmpty(t, vp.Description, "%s() variadic parameter %q has no Description", name, vp.Name)
		}
	}
}

// A cty VarParam has no upper bound, so a function faking an *optional* argument with
// one accepts any number of them and — reading only the argument it expects — drops
// the rest in silence. Both functions that do this must declare their real ceiling.
func TestExcessArgumentsAreRejected(t *testing.T) {
	tests := []struct {
		name string
		args []cty.Value
	}{
		{
			name: "bytes",
			args: []cty.Value{cty.StringVal("hi"), cty.StringVal("text/plain"), cty.StringVal("junk")},
		},
		{
			name: "base64decode",
			args: []cty.Value{cty.StringVal("aGk="), cty.StringVal("text/plain"), cty.StringVal("junk")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, ok := GetBytesFunctions()[tt.name]
			require.True(t, ok)

			_, err := fn.Call(tt.args)
			require.Error(t, err, "%s() accepted %d arguments and silently ignored the extras",
				tt.name, len(tt.args))
			assert.Contains(t, err.Error(), "at most 2 arguments")
		})
	}
}
