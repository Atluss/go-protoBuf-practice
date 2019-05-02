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
./gRPC | gRPC server for greetings requests and answers.
./protoBufApp | This directory show how to work with protoBuf protocol binary

## Greeting

Generate(in dir) with grpc plugin: `protoc --go_out=plugins=grpc:proto proto/greeting.proto`

## ProtoBufApp

* Generate proto file: `protoc --go_out=. proto/pblist.proto`  
* Install or build: `go install ./protoBufApp`, `go build ./protoBufApp`
* Run list: `protoBufApp list`
* Add list: `protoBufApp add word` where word is a string to add

#### How to read proto

* `hexdump mydb.pb`
* `hexdump -c mydb.pb`
* `cat mydb.pb`
* `cat mydb.pb | protoc --decode_raw`