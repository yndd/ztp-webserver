# syntax=docker/dockerfile:1.4
# Build the manager binary
FROM golang:1.18 as builder
WORKDIR /workspace
COPY . .
RUN --mount=type=cache,id=ztp-webserver-golang-dl-cache,target=/go/pkg/mod \
    --mount=type=cache,id=ztp-webserver-golang-build-cache,target=/root/.cache/go-build \
       CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ztp-webserver main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
#FROM gcr.io/distroless/static:nonroot
FROM alpine:latest

WORKDIR /
COPY --from=builder /workspace/ztp-webserver .
RUN mkdir -p /webserver && chown 65532:65532 /webserver
USER 65532:65532
ENTRYPOINT ["/ztp-webserver"]