linux:
	GOPRIVATE=github.com GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/overtime_linux cmd/main.go

arm:
	GOPRIVATE=github.com GOOS=linux GOARCH=arm go build -o build/overtime_arm cmd/main.go

mac:
	GOPRIVATE=github.com GOOS=darwin GOARCH=amd64 go build -o build/overtime_darwin cmd/main.go

test:
	GOPRIVATE=github.com go test -v -cover -bench . ./...

test-html:
	GOPRIVATE=github.com go test -coverprofile=coverage.out -bench . ./... && go tool cover -html=coverage.out

install:
	go build -o build/otcli cmd/main.go && sudo mv build/otcli /usr/local/bin/otcli
