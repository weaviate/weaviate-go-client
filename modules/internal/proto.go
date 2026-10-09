package internal

import proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"

func TextArray(s []string) *proto.TextArray {
	if len(s) == 0 {
		return nil
	}
	return &proto.TextArray{Values: s}
}
