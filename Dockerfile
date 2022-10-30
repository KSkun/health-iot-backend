# build
FROM golang:1.19.2-alpine3.16 AS build
COPY ./src /src
WORKDIR /src
ENV "GOPROXY" "https://goproxy.cn,direct"
RUN go build -o /build/app

# image
FROM alpine:3.16
COPY --from=build /build/app /bin/app
WORKDIR /workdir
ENTRYPOINT [ "/bin/app" ]