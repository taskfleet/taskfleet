FROM golang:1.19-buster AS build

ARG TARGETOS
ARG TARGETARCH

COPY go.* /app/
COPY pkg /app/pkg
COPY cmd /app/cmd

ENV GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOPATH=/go \
    GOCACHE=/go/cache

RUN --mount=type=cache,id=go-dependencies,target=/go \
    cd /app/cmd/gc \
    && go build -v -tags netgo -ldflags '-w -extldflags "-static"' -o app

#--------------------------------------------------------------------------------------------------

FROM debian:buster-slim

RUN apt-get update \
    && apt-get install -y ca-certificates \
    && apt-get clean
COPY --from=build /app/cmd/gc/app /app

ENTRYPOINT ["/app"]
