[![Go Report Card](https://goreportcard.com/badge/github.com/Atluss/protoBufPractice)](https://goreportcard.com/report/github.com/Atluss/protoBufPractice)

## How to use proto files

About [gRPC](https://grpc.io/docs/) and [protoBuf protocol](https://github.com/protocolbuffers/protobuf). 

After install protoc and install [plugin for go](https://github.com/golang/protobuf).
Don't forget add **~/go/bin** to your PATH, just like this for example: in **~/.profile** in your home directory add this to end of file:
```bash
if [ -d "$HOME/go/bin" ] ; then
  PATH="$PATH:$HOME/go/bin"
fi
``` 
And find it: `$PATH`, restart: `source .profile`

How to generate `*.proto` to `*.go` In folder where `*.proto` input it in terminal: `protoc --go_out=. file_name.proto`

## Project structure:

dir | description
---|---
./cmd/gRPC | gRPC server for greetings requests and answers.
./cmd/protoBufApp | This directory show how to work with protoBuf protocol binary

## Greeting

Generate(in dir) with grpc plugin: `protoc -I pkg/v1/proto/greeting --go_out=plugins=grpc:pkg/v1/proto/greeting pkg/v1/proto/greeting/greeting.proto`

## ProtoBufApp

* Generate proto file: `protoc -I pkg/v1/proto/pblist --go_out=pkg/v1/proto/pblist pkg/v1/proto/pblist/pblist.proto`  
* Install or build: `go install ./cmd/protoBufApp`, `go build ./cmd/protoBufApp`
* Run list: `protoBufApp list`
* Add list: `protoBufApp add word` where word is a string to add

#### How to read proto

* `hexdump mydb.pb`
* `hexdump -c mydb.pb`
* `cat mydb.pb`
* `cat mydb.pb | protoc --decode_raw`