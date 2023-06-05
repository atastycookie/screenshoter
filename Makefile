build_win:
	env GOOS=windows GOARCH=amd64 go build -o bin/app-amd64.exe client/app.go

build_linux:
	env GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux client/app.go

build_mac:
	env GOOS=darwin GOARCH=amd64 go build -o bin/app-amd64-darwin client/app.go

build_linux_server:
	env GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux-server server/app.go