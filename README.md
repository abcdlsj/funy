# Funy

`Funy` is a `Golang` serverless framework using `Plugin`.
Can deploy `Service` or `Function`.

> Idea from [shuttle.rs](https://www.shuttle.rs/blog/2022/04/27/dev-log-1)

~~> Currently, `Funy` dones not support `Apple Silicon`, I got a `issue` like <https://github.com/golang/go/issues/58826>.~~

## Usage

![demo.gif](demo.gif)

### Build
`Clone` the repository, then `make` to build the binary.
```fish
make
```

### Start `Daemon` `funyd`
Start `funyd` service
```fish
GIN_MODE=release USERNAME=<username> PASSWORD=<password> PORT=8080 ./bin/funyd
```

### Start `Server`
Use `funy` cli to create and deploy a server

**You need to set same `USERNAME` and `PASSWORD` in environment.**

1. create a server, with `ld_flag_x` to set `Port` and specify `main_file`
```fish
./bin/funy app --server_address=http://localhost:8080 create --name=ginserver1 --ld_flag_x="main.Port=8081" --main_file=ginserver.go
```

2. can use `curl` to get progress.

```fish
curl -su 'USERNAME:PASSWORD' localhost:8080/orchestrator/
```

output
```fish
{"ginserver1":{"process_state":0,"tar_dir":"/tmp/funy2099804210/ginserver1","main_file":"ginserver.go","ld_flag_x":{"main.Port":"8081"}}}
```

The `process_state` can be:
```go
const (
	Create ProcessState = iota
	Queued
	Built
	Loaded
	Deployed
	Error
	Deleted
)
```

3. deploy the server
```fish
./bin/funy app --server_address=http://localhost:8080 deploy --name=ginserver1 --dir=example/ginserver/
```

4. ping the server
after the `process_state` is `Deployed`, you can use `curl` to call your server.
```fish
> curl localhost:8081
{"message":"Hello World!"}⏎
> curl localhost:8081/ping
{"ping count":"1"}⏎
```

### Function
You also can deploy a function by `funy`

```fish
> ./bin/funy app --server_address=http://localhost:8080 --type=function create --name=hellohandler --main_file=handler.go
> ./bin/funy app --server_address=http://localhost:8080 --type=function deploy --name=hellohandler --dir=example/hellohandler
# When deploy success, you can use `curl` to call it
> curl -u '<username>:<password>' localhost:8080/func/hellohandler
Hello World!⏎
```

## More
**Can use `ld_flag_x` do many of things.**

I have a slightly more complex function, [example/leetjump](/example/leetjump/)

Implemented so that when you call `http://xxxx/xxx/:question_id` to jump to the issue page.
It uses `Cloudflare` `R2` to store the file containing the information underlying the `Leetcode` issue.

I declared some variables in `main.go`. You can set this using `ld_flag_x`.
```go
var (
	BUCKET          string
	ACCOUNTID       string
	ACCESSKEYID     string
	ACCESSKEYSECRET string

	LOGPREFIX string
)
```

> `LOGPREFIX` is a prefix for logs that may be used in the future.
