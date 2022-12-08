FROM --platform=$BUILDPLATFORM golang:alpine AS build
ARG BUILDPLATFORM
ARG TARGETARCH

ADD . /workspace
WORKDIR /workspace
# Build the distribution if it's not present to ensure any environment with
# docker can build, even without Go.
RUN [ -d ./dist ] || (apk add --no-cache make git && make snapshot)
RUN cp ./dist/http-wasm-tck_linux_${TARGETARCH}*/http-wasm-tck /http-wasm-tck

# Use debug distroless build instead of scratch to provide a shell for
# general troubleshooting.
FROM gcr.io/distroless/static-debian11:debug
ENTRYPOINT ["/usr/bin/http-wasm-tck"]
COPY --from=build /http-wasm-tck /usr/bin/http-wasm-tck
