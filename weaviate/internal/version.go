package internal

const clientVersion = "v5.7.5"

func GetClientVersionHeader() string {
	return "weaviate-client-go/" + clientVersion
}
