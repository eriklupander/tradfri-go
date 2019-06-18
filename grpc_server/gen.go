//go:generate protoc --go_out=paths=source_relative,plugins=grpc:golang/ tradfri.proto
package grpc_server
