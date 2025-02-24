module github.com/trashwbin/dymall/app/order

go 1.23.4

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0

replace github.com/trashwbin/dymall/rpc_gen => ../../rpc_gen

require github.com/golang/protobuf v1.5.4 // indirect
