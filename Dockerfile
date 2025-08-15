FROM golang:1.24.3-alpine3.21 AS builder
WORKDIR /resizer
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN apk add --no-cache make && make build

FROM alpine:3.21 AS runner
WORKDIR /resizer

COPY --from=builder /resizer/bin/* ./

COPY web ./web

EXPOSE 8080

CMD ["./image-resizer", "--config", "config.yaml"]