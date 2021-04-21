linux:
	GOPRIVATE=git.goasum.de GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/overtime_linux cmd/main.go

arm:
	GOPRIVATE=git.goasum.de GOOS=linux GOARCH=arm go build -o build/overtime_arm cmd/main.go

mac:
	GOPRIVATE=git.goasum.de GOOS=darwin GOARCH=amd64 go build -o build/overtime_darwin cmd/main.go

test:
	GOPRIVATE=git.goasum.de go test -v -cover -bench . ./...

test-html:
	GOPRIVATE=git.goasum.de go test -coverprofile=coverage.out -bench . ./... && go tool cover -html=coverage.out

install:
	go build -o build/otcli cmd/main.go && sudo mv build/otcli /usr/local/bin/otcli