package internal

const clientVersion = "v5.7.1"

func GetClientVersionHeader() string {
	return "weaviate-client-go/" + clientVersion
}
