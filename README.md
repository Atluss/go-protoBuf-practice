## How to use proto files

About [gRPC](https://grpc.io/docs/). 
After install protoc and install [plugin for go](https://github.com/golang/protobuf).
Don't forget add **~/go/bin** to your PATH, just like this for example: in **~/.profile** in your home directory add this to end of file:
```bash
if [ -d "$HOME/go/bin" ] ; then
  PATH="$PATH:$HOME/go/bin"
fi
``` 
And find it: `echo $PATH`, restart: `source .profile`

How to generate `*.proto` to `*.go` In folder where `*.proto` input it in terminal: `protoc --go_out=. file_name.proto`

## Greeting

Generate: `protoc --go_out=plugins=grpc:proto proto/greeting.proto`


## ProtoBuf
Install: `go install ./protoBufApp`
And run it for the test: `protoBufApp list`