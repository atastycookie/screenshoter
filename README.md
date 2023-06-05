# screenshooter

## Build
1. Install Go 1.16
2. `git clone repo`
3. `go mod vendor`
4.Change const in client/app.go

```go
const (
	downloadScript string = "http://localhost:8081/download/"
	uploadScript   string = "http://localhost:8081/upload"
)
```

5. Run make

```shell
make build_linux_server
make build_linux
make build_win
make build_mac
```

6. Для запуска бинарника под мак следует выдать права

 > System Preferences -> Security & Privacy -> Screen Recordong -> через + добавить приложение в список разрешенных.
 > Если запускается из terminal, то разрешить Teminal / iTerm
