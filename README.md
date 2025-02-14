# toll-calculator

```

```
##Installing protobuf compiler(protoc compiler)
For linux users or (wsl2)
```
sudo apt install -y protobuf-compiler
```

For Mac users can use Brew for this
```
brew install protobuff
```

##Installing GRPC and Protobuffer plugins for Golang
1. Protobuffers
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

2. GRPC
```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
3. Note that you need to set the /go/bin directory in your path
```
    PATH="${PATH}:${HOME}/go/bin"   
```

4. install package dependencies
4.1 protobuffer package
```
go get google.golang.org/protobuf
```
4.2 grpc package
```
go get google.golang.org/grpc
```

5. install prometheus
```
docker run --name prometheus -d -p 127.0.0.1:9090:9090 prom/prometheus
```
installing prom client golang
```
github.com/prometheus/client_golang/prometheus
```