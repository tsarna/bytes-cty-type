# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2026-07-21

### Added

- **`Bytes` implements the `go2cty2go` conversion interfaces**, so bytes values convert
  correctly in both directions through `go2cty2go.AnyToCty` / `CtyToAny` instead of being
  reflected over or unwrapped to a raw pointer:

  - `ToCty()` — a `*Bytes` becomes the rich bytes object (`content_type` surfaced as an
    attribute, the capsule under `_capsule`).
  - `CtyToNativeValue()` — a bytes value (bare capsule or `_capsule` object) yields its raw
    `[]byte`. The content type is dropped, since `[]byte` cannot carry it; that is the
    intended native form for serialization.

### Changed

- Now depends on `github.com/tsarna/go2cty2go` v0.3.0 (for the `CapsuleInfo` type named in
  `CtyToNativeValue`'s signature). Downstream users of this package therefore pull
  `go2cty2go` transitively.

## [0.2.0] - 2026-07-14

### Added

- **`Externs()` — the real signatures of the bytes functions, for functy hosts.**
  `externs.cty` (embedded; exposed as opaque bytes via `Externs()` and
  `ExternsFilename`) declares what all three functions actually accept, as
  [functy](https://github.com/tsarna/functy) `//functy:extern` declarations. It is never
  compiled and declares nothing callable; it exists so that `help()`, generated
  documentation, and editor tooling can show what the cty metadata cannot:

  - Each function takes an argument that is a **union** — a string, or a bytes value.
    cty has no union type, so its metadata can only say `dynamic`, which says nothing.
    Each union is declared as one form per arm.
  - `bytes` and `base64decode` each take **one optional trailing argument**, and the only
    way cty offers to make an argument optional is to make it variadic — erasing its
    name, its type, and its arity.
  - `base64decode`'s **return type depends on whether that argument is present** (string
    with one argument, bytes with two), and a cty function has exactly one signature.
    Reflected from cty alone it read as `base64decode(str, ...content_type)`, which could
    have meant anything.

  This package still does not import functy; the bytes are opaque to it. A host
  registers them:

  ```go
  parser.RegisterExterns(bytescty.Externs(), bytescty.ExternsFilename)
  ```

### Fixed

- **`bytes()` and `base64decode()` now reject excess arguments.** A cty `VarParam` has no
  upper bound of its own, so a function faking an *optional* argument with one — which is
  the only way cty offers — accepted any number of them, and its implementation, reading
  only the argument it expected, dropped the rest in silence:
  `bytes("hi", "text/plain", "junk")` returned a bytes value, and
  `base64decode("aGk=", "text/plain", "junk")` decoded happily. Both now report
  `takes at most 2 arguments`.

  This is a behavior change: a call that was silently wrong is now an error. No call that
  was doing something meaningful is affected.

### Changed

- Every function and every parameter now carries a cty `Description`. The metadata is the
  only documentation a non-functy cty host can see — and the only thing functy's own
  `doc()` reads, since `doc()` does not consult the extern.
- Depends on `rich-cty-types` v0.5.1 (was v0.4.0).

## [0.1.0] - 2026-04-17

### Added

- Initial release: the `bytes` capsule and object types, `GetBytesFunctions()`
  (`bytes`, `base64encode`, `base64decode`), and `Stringable` / `Lengthable`
  implementations for `tostring()` and `length()`.
