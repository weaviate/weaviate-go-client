package internal

const clientVersion = "5.7.0"

func GetClientVersionHeader() string {
	return "weaviate-client-go/" + clientVersion
}
