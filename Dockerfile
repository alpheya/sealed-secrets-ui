FROM golang:1.22.11-bookworm AS build
WORKDIR /go/src/app
COPY . ./

ARG VERSION

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go mod tidy \
  && go build -o /go/bin/app -ldflags="-s -w"

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/app /

CMD ["/app"]
