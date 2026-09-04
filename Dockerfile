# syntax=docker/dockerfile:1

FROM golang:1.26 AS builder

WORKDIR /src

COPY go.mod ./

COPY . .

ENV CGO_ENABLED=0

RUN go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/backend \
    ./cmd/backend

FROM gcr.io/distroless/static-debian12

COPY --from=builder /out/backend /usr/local/bin/backend

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/backend"]
