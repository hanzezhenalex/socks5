binaries: | agent cli

agent:
	go build -o $(GOPATH)/bin/agent ./src/cmd/agent/main.go

cli:
	go build -o $(GOPATH)/bin/socksctl ./src/cmd/cli/main.go

debug_agent:
	go build -gcflags="all=-N -l" -o $(GOPATH)/bin/dagent ./src/cmd/agent/main.go
	dlv --listen=:2345 --headless=true --api-version=2 exec $(GOPATH)/bin/dagent

docker_agent:
	docker build -f ./Dockerfile --target agent -t alex/socks-agent .