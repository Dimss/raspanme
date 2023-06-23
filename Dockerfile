FROM golang:1.20.1-buster as builder
ARG buildsha
WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/

RUN go build -ldflags="-X 'main.Build=${buildsha}'" -o raspan cmd/raspan/main.go

FROM ubuntu:23.04
WORKDIR /opt/raspan
COPY --from=builder /workspace/raspan raspan
COPY raspan.db /workspace/raspan.db
ENTRYPOINT ["raspan","api"]