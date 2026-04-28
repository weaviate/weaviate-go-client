package internal

const clientVersion = "v5.7.3"

func GetClientVersionHeader() string {
	return "weaviate-client-go/" + clientVersion
}
