binaries: | agent

agent:
	go build -o $(GOPATH)/bin/agent ./src/cmd/agent/main.go