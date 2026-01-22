package viper

import (
	"go-agent/gopkg/log"

	"github.com/spf13/viper"
)

type AWSConfig struct {
	AWSRegion           string `json:"aws_region" mapstructure:"aws_region"`
	AWSbucket           string `json:"aws_bucket" mapstructure:"aws_bucket"`
	AWSProxy            string `json:"aws_proxy" mapstructure:"aws_proxy"`
	AWSClientSkipVerify bool   `json:"aws_skip_verify" mapstructure:"aws_skip_verify"`
	AWSEndpointUrl      string `json:"aws_endpoint_url" mapstructure:"aws_endpoint_url"`
	AWSSecretId         string `json:"aws_secret_id" mapstructure:"aws_secret_id"`
	AWSSecretKey        string `json:"aws_secret_key" mapstructure:"aws_secret_key"`
}

func GetAws() Aws {
	logPrefix := "/gopkg/viper: viper.GetAws()"
	// 解析ES配置
	var cfgStruct AWSConfig
	if err := viper.UnmarshalKey("aws_s3", &cfgStruct); err != nil {
		log.Sugar().Error(logPrefix, "Failed to parse aws_s3 error", err.Error())
	}

	return Aws{
		Region: cfgStruct.AWSRegion,
		Bucket: cfgStruct.AWSbucket,
		Proxy:  cfgStruct.AWSProxy,
		Client: AwsClient{
			SkipVerify: cfgStruct.AWSClientSkipVerify,
		},
		Endpoint: AwsEndpoint{
			URL: cfgStruct.AWSEndpointUrl,
		},
		Secret: AwsSecret{
			ID:  cfgStruct.AWSSecretId,
			Key: cfgStruct.AWSSecretKey,
		},
	}
}

type MinioConfig struct {
	Endpoint  string `json:"endpoint" mapstructure:"endpoint"`
	SecretId  string `json:"secret_id" mapstructure:"secret_id"`
	SecretKey string `json:"secret_key" mapstructure:"secret_key"`
	UseSSL    bool   `json:"use_ssl" mapstructure:"use_ssl"`
	Bucket    string `json:"bucket" mapstructure:"bucket"`
}

func GetMinioCnf() Minio {
	logPrefix := "/gopkg/viper: viper.GetAws()"
	// 解析ES配置
	var cfgStruct MinioConfig
	if err := viper.UnmarshalKey("minio", &cfgStruct); err != nil {
		log.Sugar().Error(logPrefix, "Failed to parse minio error", err.Error())
	}

	return Minio{
		Endpoint:  cfgStruct.Endpoint,
		Bucket:    cfgStruct.Bucket,
		SecretId:  cfgStruct.SecretId,
		SecretKey: cfgStruct.SecretKey,
		UseSSL:    cfgStruct.UseSSL,
	}
}
