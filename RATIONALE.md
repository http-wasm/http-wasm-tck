# Notable rationale of http-wasm-tck runner

## Using TestMain

This project is a standalone binary but does not use the typical pattern of a
package named `main` with a function named `main`. Instead, it uses a test file,
`runner_test.go`, which implements `TestMain` with functionality, not all
related to running tests. This provides a `testing.T` which can be passed to
the TCK, so it can use the Go test framework, for example to provide output of
each test case.

An alternative that would use the more idiomatic `main` package could
be used by setting up the testing framework manually. However, Go does not
provide a clear mechanism for doing this. [testing.Main][1] is documented as
being legacy and systems should be updated to `testing.MainStart`.
[testing.MainStart][2] is documented as not meaning to be called directly.

Implementing `TestMain` and compiling our binary with `go test -c` gives us
the functionality we need without relying on any internal or questionably
supported mechanisms.

## Development

The entrypoint to this application is in [runner_test.go][3], not a `main`
package as is typical. See the [rationale][4] for more information on why.

[1]: https://pkg.go.dev/testing#Main
[1]: https://pkg.go.dev/testing#MainStart
[3]: ./runner_test.go
[4]: ./RATIONALE.md
