FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build

FROM alpine:3.18

WORKDIR /app

COPY --from=build /app/genesis /app

EXPOSE 8080

CMD ["./genesis"]
