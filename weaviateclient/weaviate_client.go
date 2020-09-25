package weaviateclient

const apiVersion = "v1"

type Config struct {
	Host string
	Scheme string
}

func (cfg *Config) basePath() string {
	return cfg.Scheme+"://"+cfg.Host+"/"+apiVersion

}

type WeaviateClient struct {
	config Config
	Misc Misc
}

func New(config Config) *WeaviateClient {

	//swagger.Configuration{
	//	BasePath:      "",
	//	Host:          "",
	//	Scheme:        "",
	//	DefaultHeader: nil,
	//	UserAgent:     "",
	//	HTTPClient:    nil,
	//}
	//cfg := &swagger.Configuration{
	//	BasePath:      "https://localhost/v1",
	//	DefaultHeader: make(map[string]string),
	//	UserAgent:     "Swagger-Codegen/1.0.0/go",
	//}

	return &WeaviateClient{
		config: config,
		Misc: Misc{config: &config},
	}
}








