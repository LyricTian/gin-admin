FROM golang:alpine as builder

ARG APP=ginadmin
ARG VERSION=v10.1.0
ARG RELEASE_TAG=$(VERSION)

# Install the required packages
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set CGO_CFLAGS to enable large file support
ENV CGO_CFLAGS "-D_LARGEFILE64_SOURCE"

ENV GOPROXY="https://goproxy.cn"

WORKDIR /go/src/${APP}
COPY . .

# Build the application
RUN go build -ldflags "-w -s -X main.VERSION=${RELEASE_TAG}" -o ./${APP} .

FROM alpine
ARG APP=ginadmin
WORKDIR /go/src/${APP}
COPY --from=builder /go/src/${APP}/${APP} /usr/bin/
# COPY --from=builder /go/src/${APP}/configs /usr/bin/configs
# COPY --from=builder /go/src/${APP}/dist /usr/bin/dist
ENTRYPOINT ["ginadmin", "start", "-d", "/usr/bin/configs", "-c", "prod", "-s", "/usr/bin/dist"]
EXPOSE 8040
