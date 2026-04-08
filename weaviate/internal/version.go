package internal

const clientVersion = "v5.7.2"

func GetClientVersionHeader() string {
	return "weaviate-client-go/" + clientVersion
}
