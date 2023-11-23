# Funy

`Funy` is a `Golang` serverless framework using `Plugin`.
> Idea from [shuttle.rs](https://www.shuttle.rs/blog/2022/04/27/dev-log-1)

> Currently, `Funy` dones not support `Apple Silicon`, I got a `issue` like <https://github.com/golang/go/issues/58826>.

## Usage

1. `Clone` the repository, then `make` to build the binary.
```fish
make
```

2. Start `funyd` service
```fish
./bin/funyd
```

3. Use `funy` cli to create and deploy a server
```fish
funy (main)> ls
README.md  cmd/  example/  go.mod  go.sum  internal/  makefile
funy (main)> ./bin/funy app --server_address=http://localhost:8080 create --name=ginserver1 --ld_flag_x="main.Port=8081" --main_file=ginserver.go
funy (main)> curl localhost:8080
{"ginserver1":{"process_state":0,"tar_dir":"/tmp/funy2099804210/ginserver1","main_file":"ginserver.go","ld_flag_x":{"main.Port":"8081"}}}⏎
funy (main)> ./bin/funy app --server_address=http://localhost:8080 deploy --name=ginserver1 --dir=example/ginserver/
funy (main)> curl localhost:8080
{"ginserver1":{"process_state":1,"tar_dir":"/tmp/funy2099804210/ginserver1","main_file":"ginserver.go","ld_flag_x":{"main.Port":"8081"}}}⏎
funy (main)> curl localhost:8080
{"ginserver1":{"process_state":4,"tar_dir":"/tmp/funy2099804210/ginserver1","main_file":"ginserver.go","ld_flag_x":{"main.Port":"8081"}}}⏎
funy (main)> curl localhost:8081
{"message":"Hello World!"}⏎
funy (main)> curl localhost:8081/ping
{"ping count":"1"}⏎
```