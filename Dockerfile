FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o field-service .

FROM alpine:latest

RUN apk --no-cache add tzdata

WORKDIR /app

EXPOSE 8003

COPY --from=builder /app/field-service /app/field-service

ENTRYPOINT ["/app/field-service"]
