# build stage
FROM golang:1.12.1-alpine AS build

RUN apk add --no-cache git dep

WORKDIR /go/src/github.com/sdorra/host-path-provisioner
COPY Gopkg.toml .
COPY Gopkg.lock .

RUN dep ensure -vendor-only

COPY . .

RUN go build -o host-path-provisioner

# final stage
FROM alpine:3.9.2
COPY --from=build /go/src/github.com/sdorra/host-path-provisioner/host-path-provisioner /bin/host-path-provisioner

VOLUME /volumes

ENTRYPOINT /bin/host-path-provisioner
CMD ["-directory /volumes"]
