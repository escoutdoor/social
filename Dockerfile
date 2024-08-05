FROM golang:1-alpine3.20 as builder
LABEL authors="escoutdoor"

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# RUN go build -o ./bin/social ./cmd/api/main.go
#
# FROM alpine:3.20
# WORKDIR /app
# COPY --from=builder /app/bin/social /app/

CMD ["air"]
