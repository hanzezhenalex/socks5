binaries: | agent

agent:
	go build -o $(GOPATH)/bin/vflow ./src/cmd/agent/main.go