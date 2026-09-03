FROM --platform=$BUILDPLATFORM golang:1.22.11-bookworm AS build
WORKDIR /go/src/app
COPY . ./

ARG VERSION
# Set by buildx per target platform (e.g. "amd64", "arm64") when building
# with --platform linux/amd64,linux/arm64 -- upstream hardcoded GOARCH=amd64
# here, which is exactly why the published 0.3.3 image only ever ran on
# amd64 nodes and crashed with "exec format error" on arm64 (e.g. Oracle
# Ampere A1). Cross-compiling with $TARGETARCH instead of hardcoding fixes
# that for any platform buildx is asked to target, not just arm64.
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

RUN go mod tidy \
  && go build -o /go/bin/app -ldflags="-s -w"

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/app /

CMD ["/app"]
