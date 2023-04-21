# http-wasm technology compatibility kit (TCK) runner

The TCK is a test suite for checking conformance of http-wasm implementations
versus the http-wasm [ABI specification][1].

This repository contains a standalone runner executable for the
[http-wasm TCK][2] implemented in Go for use with non-Go host implementations.

## Running the TCK

The basic steps for running the TCK are

1. Implement the backend handler, which is the business logic wrapped by
middleware
2. Set up the middleware using the TCK guest wasm module
3. Start an HTTP server serving this middleware
4. Run the tests, pointing at the URL for the HTTP server

Precompiled executables can be downloaded from the [releases][3].
A Docker image is also provided as [ghcr.io/http-wasm/http-wasm-tck][4].

The Go implementation of the backend handler can be referred to [here][5].

An HTTP server using this implementation can be started using
`http-wasm-tck backend`. The server defaults to listening on
`0.0.0.0:9080`, which can be changed with the `-addr` flag.
If the port passed to address is `0`, for example with
`http-wasm-tck backend -addr 0.0.0.0:0`, a random port will be chosen. The
backend always prints a message like `Started backend server on 0.0.0.0:9080`
which can be parsed to find the port if random.

Alternatively, it may be simpler to reimplement the Go implementation of the
backend handler into the language of the host.

The guest wasm module can be extracted from the command using
`http-wasm-tck extract-guest <path to wasm file>`.

With the HTTP server started and serving the middleware and backend, the tests
can be run using `http-wasm-tck run`. This defaults to issuing requests to
`http://localhost:8080`, which can be changed with the `-url` flag.

## Development

The entrypoint to this application is in [runner_test.go][3], not a `main`
package as is typical. See the [rationale][6] for more information on why.

[1]: https://http-wasm.io/http-handler-abi/
[2]: https://github.com/http-wasm/http-wasm-host-go/tree/main/tck
[3]: https://github.com/http-wasm/http-wasm-tck/releases
[4]: https://github.com/http-wasm/http-wasm-tck/pkgs/container/http-wasm-tck
[5]: https://github.com/http-wasm/http-wasm-host-go/blob/359f2659391c4407272406a818dfc8bdef934419/tck/backend.go#L13
[6]: ./RATIONALE.md
