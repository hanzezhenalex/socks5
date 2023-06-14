FROM golang:1.19 as build

# Create app directory
WORKDIR /usr/src/app

COPY . .

RUN make binaries

FROM busybox as agent

ENTRYPOINT ["${GOPATH}/bin/agent", "start"]