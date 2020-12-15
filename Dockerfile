FROM golang:1.15.6 as builder

WORKDIR /build
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/inventory cmd/inventory/main.go

FROM busybox:1.32.0

WORKDIR /app
COPY --from=builder /build/bin/inventory .

ENTRYPOINT ["/app/inventory"]
