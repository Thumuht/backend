## build
FROM golang:1.19-alpine AS build
WORKDIR /app
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
COPY . .
RUN go build -o / ./...

## deploy
FROM alpine:latest
WORKDIR /
COPY --from=build /cmd /cmd
ENTRYPOINT [ "/cmd" ]