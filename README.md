# spaas-gem challenge

### Compile and Run

to compile the API, use the following commands:

```bash
$ cd <project-dir>
$ go mod download
$ go build -o gem cmd/gem/main.go
```

note: this application has been written with `go 1.17`

to run the application (from the directory where the application has been built):

```bash
$ ./gem
# to specify a specific listenning port, use:
$ ./gem -a ip:port
```

the binary accepts also a `[-a]` option to specify the address where the embedded
HTTP server will listen for incoming requests. The default is to listen on all
interface on port `8888`

the API, as requested, expose a `/productionplan` endpoint.

to subscribe to the websocket socket service, you can submit a request to the `/ws`
endpoint.

### Test

to execute the test, you can execute in the project directory:

```bash
$ go test -v
```
