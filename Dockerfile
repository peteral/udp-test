FROM golang:1.21-alpine as build

WORKDIR /opt

COPY ./ ./

RUN go build -o bin/udp-test .

# stage 2 - build application image
FROM alpine
RUN apk --no-cache add ca-certificates

ENV GODEBUG=netdns=go

WORKDIR /opt

COPY --from=build /opt/bin/udp-test ./bin/udp-test

EXPOSE 80
ENTRYPOINT ["./bin/udp-test"]
