package bytescty

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

// --- bytes() ---

func TestBytesFunc_FromString(t *testing.T) {
	fn := MakeBytesFunc()

	result, err := fn.Call([]cty.Value{cty.StringVal("hello")})
	require.NoError(t, err)
	assert.Equal(t, BytesObjectType, result.Type())
	b, err := GetBytesFromValue(result)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), b.Data)
	assert.Equal(t, "", b.ContentType)
}

func TestBytesFunc_FromStringWithContentType(t *testing.T) {
	fn := MakeBytesFunc()

	result, err := fn.Call([]cty.Value{cty.StringVal("hello"), cty.StringVal("text/plain")})
	require.NoError(t, err)
	assert.Equal(t, BytesObjectType, result.Type())
	b, err := GetBytesFromValue(result)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), b.Data)
	assert.Equal(t, "text/plain", b.ContentType)
}

func TestBytesFunc_FromBytesObject_PreservesContentType(t *testing.T) {
	fn := MakeBytesFunc()
	src := BuildBytesObject([]byte{1, 2, 3}, "image/png")

	result, err := fn.Call([]cty.Value{src})
	require.NoError(t, err)
	b, err := GetBytesFromValue(result)
	require.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, b.Data)
	assert.Equal(t, "image/png", b.ContentType)
}

func TestBytesFunc_FromBytesObject_OverridesContentType(t *testing.T) {
	fn := MakeBytesFunc()
	src := BuildBytesObject([]byte{1, 2, 3}, "image/png")

	result, err := fn.Call([]cty.Value{src, cty.StringVal("image/jpeg")})
	require.NoError(t, err)
	b, err := GetBytesFromValue(result)
	require.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, b.Data)
	assert.Equal(t, "image/jpeg", b.ContentType)
}

func TestBytesFunc_InvalidType(t *testing.T) {
	fn := MakeBytesFunc()
	_, err := fn.Call([]cty.Value{cty.NumberIntVal(42)})
	assert.Error(t, err)
}

// --- base64encode() ---

func TestBase64EncodeFunc_FromString(t *testing.T) {
	fn := MakeBase64EncodeFunc()
	result, err := fn.Call([]cty.Value{cty.StringVal("hello")})
	require.NoError(t, err)
	assert.True(t, result.RawEquals(cty.StringVal("aGVsbG8=")))
}

func TestBase64EncodeFunc_FromBytesObject(t *testing.T) {
	fn := MakeBase64EncodeFunc()
	b := BuildBytesObject([]byte("hello"), "")
	result, err := fn.Call([]cty.Value{b})
	require.NoError(t, err)
	assert.True(t, result.RawEquals(cty.StringVal("aGVsbG8=")))
}

func TestBase64EncodeFunc_InvalidType(t *testing.T) {
	fn := MakeBase64EncodeFunc()
	_, err := fn.Call([]cty.Value{cty.NumberIntVal(1)})
	assert.Error(t, err)
}

// --- base64decode() ---

func TestBase64DecodeFunc_OneArg_ReturnsString(t *testing.T) {
	fn := MakeBase64DecodeFunc()
	result, err := fn.Call([]cty.Value{cty.StringVal("aGVsbG8=")})
	require.NoError(t, err)
	assert.Equal(t, cty.String, result.Type())
	assert.True(t, result.RawEquals(cty.StringVal("hello")))
}

func TestBase64DecodeFunc_TwoArgs_ReturnsBytesWithContentType(t *testing.T) {
	fn := MakeBase64DecodeFunc()
	result, err := fn.Call([]cty.Value{cty.StringVal("aGVsbG8="), cty.StringVal("text/plain")})
	require.NoError(t, err)
	assert.Equal(t, BytesObjectType, result.Type())
	b, err := GetBytesFromValue(result)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), b.Data)
	assert.Equal(t, "text/plain", b.ContentType)
}

func TestBase64DecodeFunc_TwoArgs_EmptyContentType(t *testing.T) {
	fn := MakeBase64DecodeFunc()
	result, err := fn.Call([]cty.Value{cty.StringVal("aGVsbG8="), cty.StringVal("")})
	require.NoError(t, err)
	assert.Equal(t, BytesObjectType, result.Type())
	b, err := GetBytesFromValue(result)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), b.Data)
	assert.Equal(t, "", b.ContentType)
}

func TestBase64DecodeFunc_InvalidBase64(t *testing.T) {
	fn := MakeBase64DecodeFunc()
	_, err := fn.Call([]cty.Value{cty.StringVal("not-valid-base64!!!")})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "base64decode")
}
