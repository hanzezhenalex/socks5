binaries: | agent

agent:
	go build -o $(GOPATH)/bin/agent ./src/cmd/agent/main.go

docker_agent:
	docker build -f ./Dockerfile --target agent -t alex/socks-agent .