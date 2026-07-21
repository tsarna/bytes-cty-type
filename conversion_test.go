package bytescty

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsarna/go2cty2go"
)

func TestToCty_ProducesRichObject(t *testing.T) {
	// A *Bytes flowing through go2cty2go.AnyToCty becomes the rich object,
	// not a reflected struct.
	b := &Bytes{Data: []byte("hi"), ContentType: "text/plain"}
	got, err := go2cty2go.AnyToCty(b)
	require.NoError(t, err)

	require.True(t, got.Type().IsObjectType(), "want the rich bytes object, got %s", got.Type().GoString())
	assert.Equal(t, "text/plain", got.GetAttr("content_type").AsString())
	assert.True(t, got.Type().HasAttribute("_capsule"))
}

func TestCtyToNativeValue_BareCapsuleYieldsBytes(t *testing.T) {
	cap := NewBytesCapsule([]byte("raw"), "application/octet-stream")
	got, err := go2cty2go.CtyToAny(cap)
	require.NoError(t, err)
	assert.Equal(t, []byte("raw"), got, "a bytes capsule must yield []byte, not the *Bytes pointer")
}

func TestCtyToNativeValue_RichObjectYieldsBytes(t *testing.T) {
	obj := BuildBytesObject([]byte("raw"), "text/plain")
	got, err := go2cty2go.CtyToAny(obj)
	require.NoError(t, err)
	assert.Equal(t, []byte("raw"), got,
		"a bytes object must yield []byte, not a map carrying the raw capsule pointer")
}

func TestCtyToNativeValue_BinaryPreserved(t *testing.T) {
	raw := []byte{0xff, 0x00, 0xfe, 0x80}
	got, err := go2cty2go.CtyToAny(BuildBytesObject(raw, ""))
	require.NoError(t, err)
	assert.Equal(t, raw, got)
}

func TestRoundTrip_BytesThroughGo2Cty(t *testing.T) {
	// AnyToCty(*Bytes) -> rich object -> CtyToAny -> []byte. Not identity
	// (native []byte cannot carry content_type), but the payload survives.
	b := &Bytes{Data: []byte("payload"), ContentType: "text/plain"}
	cv, err := go2cty2go.AnyToCty(b)
	require.NoError(t, err)
	back, err := go2cty2go.CtyToAny(cv)
	require.NoError(t, err)
	assert.Equal(t, []byte("payload"), back)
}
