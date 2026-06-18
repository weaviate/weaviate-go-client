package internal

const clientVersion = "v5.7.4"

func GetClientVersionHeader() string {
	return "weaviate-client-go/" + clientVersion
}
