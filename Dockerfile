FROM golang:1.19 as build

# Create app directory
WORKDIR /usr/src/app

COPY . .

RUN make binaries

# busybox
FROM golang:1.19 as agent

COPY --from=build /go/bin/agent /usr/bin/agent
COPY --from=build /go/bin/socksctl /usr/bin/socksctl

ENTRYPOINT ["/usr/bin/agent"]