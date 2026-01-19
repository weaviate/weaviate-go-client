package proto

//go:generate sh -c "rm -rf v1 && mkdir v1 && ../../../../../bin/protoc --proto_path=../../../../../api/proto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ../../../../../api/proto/v1/*.proto"
