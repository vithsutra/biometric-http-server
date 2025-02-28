FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/main ./cmd/*

FROM alpine:latest

WORKDIR /app

RUN mkdir -p font-family assets

COPY --from=build /app/font-family ./font-family

COPY --from=build /app/assets ./assets

COPY --from=build /app/bin/main .

ENTRYPOINT [ "./main" ]