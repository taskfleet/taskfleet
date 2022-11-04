FROM golang:1.19-buster AS build

ARG TARGETOS
ARG TARGETARCH

RUN go install github.com/grpc-ecosystem/grpc-health-probe@latest

WORKDIR /app

COPY go.mod go.sum \
    ./
COPY packages/dymant packages/eagle packages/mercury \
    packages/
COPY services/instance-manager/cmd services/instance-manager/internal \
    services/instance-manager/

ENV GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOPATH=/go \
    GOCACHE=/go/cache

RUN --mount=type=cache,id=go-dependencies,target=/go \
    cd services/instance-manager/cmd/service \
    && go build -v -o app

#--------------------------------------------------------------------------------------------------

FROM debian:buster

RUN apt-get update \
    && apt-get install -y ca-certificates \
    && apt-get clean
COPY --from=build /app/services/instance-manager/cmd/service/app /app
COPY --from=build /go/bin/grpc-health-probe /bin/grpc-health-probe

ENTRYPOINT ["/app"]
EXPOSE 5404
