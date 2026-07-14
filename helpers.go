package bytescty

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// boundedArity returns a Type func that rejects too many arguments and then reports
// a fixed return type.
//
// It exists because a cty VarParam has no upper bound of its own. A function that
// fakes an *optional* argument with a variadic — which is the only way cty offers —
// therefore accepts any number of trailing arguments, and its Impl, reading only the
// one it expects, drops the rest in silence: bytes("hi", "text/plain", "junk") used to
// return a bytes value. Declaring the ceiling is how such a function says how many
// arguments it actually takes.
func boundedArity(name string, maxArgs int, ret cty.Type) function.TypeFunc {
	return func(args []cty.Value) (cty.Type, error) {
		if len(args) > maxArgs {
			return cty.NilType, fmt.Errorf("%s() takes at most %d arguments, got %d", name, maxArgs, len(args))
		}
		return ret, nil
	}
}
