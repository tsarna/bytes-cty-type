# bytes-cty-type

A [go-cty](https://github.com/zclconf/go-cty) capsule type for immutable byte slices with optional MIME content types, plus functions for encoding and decoding.

[![CI](https://github.com/tsarna/bytes-cty-type/actions/workflows/ci.yml/badge.svg)](https://github.com/tsarna/bytes-cty-type/actions/workflows/ci.yml)

## Overview

This package provides a `bytes` type for use in HCL2 / cty expression evaluation contexts. A bytes value carries raw binary data and an optional content type string. It is represented as a cty object with a `content_type` attribute and an internal `_capsule` attribute for interface dispatch.

Extracted from [vinculum](https://github.com/tsarna/vinculum), where it powers binary payloads in HTTP, serialization, and key-value store operations.

## Types

### `bytescty.Bytes`

```go
type Bytes struct {
    Data        []byte
    ContentType string
}
```

### `bytescty.BytesCapsuleType`

A cty capsule type wrapping `*Bytes`. Used internally as the `_capsule` attribute of bytes objects.

### `bytescty.BytesObjectType`

A cty object type with attributes:
- `content_type` (string) — the MIME type
- `_capsule` (BytesCapsuleType) — the encapsulated bytes

### Helper functions

```go
bytescty.NewBytesCapsule(data []byte, contentType string) cty.Value
bytescty.BuildBytesObject(data []byte, contentType string) cty.Value
bytescty.GetBytesFromCapsule(val cty.Value) (*Bytes, error)
bytescty.GetBytesFromValue(val cty.Value) (*Bytes, error)
```

`GetBytesFromValue` accepts a raw capsule, a bytes object, or any cty object with a `_capsule` attribute containing a `*Bytes` — it delegates to [rich-cty-types](https://github.com/tsarna/rich-cty-types) `GetCapsuleFromValue`.

## Registration

```go
import bytescty "github.com/tsarna/bytes-cty-type"

// Add all bytes functions to your eval context:
for name, fn := range bytescty.GetBytesFunctions() {
    funcs[name] = fn
}
```

### rich-cty-types integration

The `Bytes` type implements the [rich-cty-types](https://github.com/tsarna/rich-cty-types) `Stringable` and `Lengthable` interfaces:

- `tostring(b)` returns the raw bytes as a UTF-8 string.
- `length(b)` returns the byte length.

```go
import (
    bytescty "github.com/tsarna/bytes-cty-type"
    richcty  "github.com/tsarna/rich-cty-types"
)

funcs := richcty.GetGenericFunctions()           // tostring, length, ...
for name, fn := range bytescty.GetBytesFunctions() {
    funcs[name] = fn
}
```

## Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `bytes(s)` | `(string) → bytes` | Create bytes from a UTF-8 string |
| `bytes(s, ct)` | `(string, string) → bytes` | Create bytes with a content type |
| `bytes(b)` | `(bytes) → bytes` | Copy a bytes value (preserves content type) |
| `bytes(b, ct)` | `(bytes, string) → bytes` | Copy with overridden content type |
| `base64encode(v)` | `(string\|bytes) → string` | Encode a string or bytes value to base64 |
| `base64decode(s)` | `(string) → string` | Decode base64 to string (backward compatible) |
| `base64decode(s, ct)` | `(string, string) → bytes` | Decode base64 to bytes object with content type |

## Examples

```hcl
# Create bytes from a string
bytes("hello world")
bytes("hello", "text/plain")

# Access content type
bytes("img data", "image/png").content_type   # "image/png"

# Base64 round-trip
base64encode("hello")                          # "aGVsbG8="
base64decode("aGVsbG8=")                       # "hello" (string)
base64decode("aGVsbG8=", "text/plain")         # bytes object

# Re-wrap with different content type
bytes(existing_bytes, "application/json")

# With rich-cty-types generic functions
tostring(b)                                    # raw bytes as string
length(b)                                      # byte count
```

## License

BSD 2-Clause — see [LICENSE](LICENSE).
