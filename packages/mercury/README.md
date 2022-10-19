# Go Package Mercury

Mercury provides a common interface for deploying applications running inside Kubernetes. The core
concept is a _runtime_ to which different long-running services can be attached. Whenever one of
the services fails, the entire runtime fails.

Currently, Mercury provides support for the following two services:

- gRPC servers with health checks
- HTTP servers for Prometheus metrics

## Installation

```bash
go get go.taskfleet.io/packages/mercury
```
