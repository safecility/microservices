### Milesight Protobuffer
// install protoc and add go support:
````
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

protoc --go_out=./  ./usage.proto

go install github.com/GoogleCloudPlatform/protoc-gen-bq-schema@latest

````


