FROM golang:1.22.1-alpine3.19 as builder

ENV GOPROXY=https://goproxy.io,direct
RUN apk update --no-cache
WORKDIR /app
COPY . /app
RUN go clean --modcache
RUN go build -mod=readonly -o app cmd/app.go

FROM alpine

RUN apk update --no-cache
WORKDIR /app
COPY --from=builder /app /app

CMD ./app