package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TykTechnologies/tyk/ctx"
	"github.com/TykTechnologies/tyk/log"
	"github.com/TykTechnologies/tyk/user"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"net/http"
)

func InvokeLambda(rw http.ResponseWriter, r *http.Request) {
	// Read Static AWS Configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"AWS_ACCESS_KEY_ID",
			"AWS_SECRET_ACCESS_KEY",
			"AWS_SESSION_TOKEN")),
	)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}

	// If the lambda function to be invoked is not specified in the header, execute the default
	function := r.Header.Get("function")
	if len(function) == 0 {
		function = "jia_test_custom_go_plugin"
	}

	lambdaClient := lambda.NewFromConfig(cfg)
	payload, err := json.Marshal("")
	invokeOutput, err := lambdaClient.Invoke(context.TODO(), &lambda.InvokeInput{
		FunctionName: aws.String(function),
		LogType:      types.LogTypeNone,
		Payload:      payload,
	})

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(invokeOutput.Payload[:])
	if err != nil {
		return
	}

}

func main() {}

func init() {
	logger.Info("--- Go custom plugin v4 init success! ---- ")
}
