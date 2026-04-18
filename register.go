package bytescty

import "github.com/zclconf/go-cty/cty/function"

// GetBytesFunctions returns bytes-related cty functions for registration in an eval context.
func GetBytesFunctions() map[string]function.Function {
	return map[string]function.Function{
		"bytes":        MakeBytesFunc(),
		"base64encode": MakeBase64EncodeFunc(),
		"base64decode": MakeBase64DecodeFunc(),
	}
}
