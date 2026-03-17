//go:generate protoc --go_out=golang/ --go_opt=paths=source_relative --go-grpc_out=golang/ --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative tradfri.proto
package grpc_server
