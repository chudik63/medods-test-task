FROM golang:1.23 AS builder

WORKDIR /build

COPY go.mod go.sum ./     

RUN go clean --modcache && go mod download      

COPY . .       

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/main .
COPY --from=builder /build/.env .
COPY --from=builder /build/docs ./docs

CMD ["./main"]