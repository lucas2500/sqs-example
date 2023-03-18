package util

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	Cfg aws.Config
)

func LoadAWSConfig() {

	var err error

	Cfg, err = config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatal("- Unable to load default config")
	}
}
