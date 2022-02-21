FROM golang:1.17.7 AS build-env

ARG VERSION
ARG WORKDIR
ARG PORT
ARG ENV

ADD . $WORKDIR
WORKDIR $WORKDIR

RUN go mod download
RUN CGO_ENABLED=0 go build \
    -o main \
    -mod=readonly \
    -ldflags "-extldflags 'static' -X main.version=${VERSION}" main.go

FROM alpine as cert
RUN apk update && apk add ca-certificates

FROM alpine:latest
ARG WORKDIR
ARG ENV

RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/

COPY --from=cert /etc/ssl/certs /etc/ssl/certs
COPY --from=build-env $WORKDIR/main /usr/local/bin/server
COPY --from=build-env $WORKDIR/configs /usr/local/bin/configs

EXPOSE $PORT
CMD ENV=$ENV /usr/local/bin/server
