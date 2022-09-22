FROM golang:alpine as builder

ARG APP=ginadmin
ARG VERSION=v9.0.1
ARG RELEASE_TAG=$(VERSION)

WORKDIR /go/src/${APP}
COPY . .
RUN go build -ldflags "-w -s -X main.VERSION=${RELEASE_TAG}" -o ./${APP} .

FROM alpine
ARG APP=ginadmin
WORKDIR /go/src/${APP}
COPY --from=builder /go/src/${APP}/${APP} /usr/bin/
COPY --from=builder /go/src/${APP}/configs /usr/bin/configs
ENTRYPOINT ["ginadmin", "start", "--configdir", "/usr/bin/configs"]
EXPOSE 8080