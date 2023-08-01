FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build

FROM alpine:3.18

WORKDIR /app

COPY --from=build /app/genesis /app

EXPOSE 8080

CMD ["./genesis"]
