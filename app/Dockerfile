FROM golang:1.23.4-alpine AS builder

LABEL stage=gobuilder
LABEL authors="mos_gor"

ENV CGO_ENABLED 0

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o /app/main cmd/CalorieGuide-db/main.go


FROM alpine

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /build/config ./config

CMD ["./main"]