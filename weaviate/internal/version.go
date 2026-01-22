package internal

const clientVersion = "5.6.0"

func GetClientVersionHeader() string {
	return "weaviate-client-go/" + clientVersion
}
