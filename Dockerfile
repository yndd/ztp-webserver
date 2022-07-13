# syntax=docker/dockerfile:1.3
# Build the manager binary
FROM golang:1.17 as builder
WORKDIR /workspace
ENV GOPATH /
# Copy the Go Modules manifests
COPY go.mod go.sum ./
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN --mount=type=cache,id=ztp-webserver-golang-cache,target=/go/pkg/mod go mod download
# Copy the go source
COPY . .
# Build
RUN --mount=type=cache,id=ztp-webserver-golang-cache,target=/go/pkg/mod CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ztp-webserver main.go


# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
#FROM gcr.io/distroless/static:nonroot
FROM alpine:latest

WORKDIR /
COPY --from=builder /workspace/ztp-webserver .
RUN mkdir -p /webserver && chown 65532:65532 /webserver
USER 65532:65532
ENTRYPOINT ["/ztp-webserver"]