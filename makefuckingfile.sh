SMC_AGENT_PORT=9898
SMC_SHOW_IMAGE=true
SMC_SHOW_NHSO=true
SMC_SHOW_LASER=true
export GOOS=windows
export GOARCH=amd64
export SMC_PORT
export SMC_SHOW_IMAGE
export SMC_SHOW_NHSO
export SMC_SHOW_LASER
go build -o ./bin/thai-smartcard-agent-windows-amd64.exe ./cmd/agent/main.go
go build -ldflags "-H windowsgui" -o ./bin/thai-smartcard-agent.windows-amd64-no-console.exe ./cmd/agent/main.go
