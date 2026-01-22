package viper

type Aws struct {
	Region   string      `json:"region"`
	Bucket   string      `json:"bucket"`
	Proxy    string      `json:"proxy"`
	Endpoint AwsEndpoint `json:"endpoint"`
	Secret   AwsSecret   `json:"secret"`
	Client   AwsClient   `json:"client"`
}

type AwsClient struct {
	SkipVerify bool `json:"skip_verify"`
}

type AwsEndpoint struct {
	URL string `json:"url"`
}

type AwsSecret struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

type Minio struct {
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	UseSSL    bool   `json:"use_ssl"`
}
