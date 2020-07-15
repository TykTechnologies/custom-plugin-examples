package main

import (
	"context"
	"log"
	"os"

	coprocess "github.com/asoorm/tyk-mw-grpcgo-lambda/proto/go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
)

type Dispatcher struct {
	lambda *lambda.Lambda
}

func NewDispatcher() *Dispatcher {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := lambda.New(sess, &aws.Config{Region: aws.String(awsRegion)})

	lfo, err := svc.ListFunctions(nil)
	if err != nil {
		log.Println(errors.Wrap(err, "unable to list avail functions"))
		os.Exit(1)
		return nil
	}

	log.Println("### Available Functions ###")
	for _, f := range lfo.Functions {
		log.Println("---> ", *f.FunctionName)
	}
	log.Println("###")

	log.Println("in your api definition, configure the custom middleware as follows, replacing FUNC_NAME with one of available functions printed above:")
	log.Println(`"custom_middleware":{"pre":[{"name":"FUNC_NAME"}],"post":[],"post_key_auth":[],"auth_check":{"name":"","path":"","require_session":false},"response":[],"driver":"grpc","id_extractor":{"extract_from":"","extract_with":"","extractor_config":{}}}`)

	return &Dispatcher{
		lambda: svc,
	}
}

func (d Dispatcher) Dispatch(ctx context.Context, obj *coprocess.Object) (*coprocess.Object, error) {
	log.Println("Dispatch called")

	switch obj.HookName {
	case "":
		log.Println("dont know what to do with empty hook")
	default:
		log.Printf("invoking %s function", obj.HookName)

		out, err := d.lambda.Invoke(&lambda.InvokeInput{
			FunctionName:   aws.String(obj.HookName),
			Payload:        []byte(`{"foo":"bar"}`),
			InvocationType: aws.String("RequestResponse")},
		)
		if err != nil {
			log.Println(errors.Wrap(err, "error invoking function"))
			obj.Request.ReturnOverrides.ResponseCode = 500
			obj.Request.ReturnOverrides.ResponseError = string("error invoking lambda function")

			return obj, nil
		}

		obj.Request.ReturnOverrides.ResponseCode = 200
		obj.Request.ReturnOverrides.ResponseError = string(out.Payload)
	}

	return obj, nil
}

func (d Dispatcher) DispatchEvent(ctx context.Context, obj *coprocess.Event) (*coprocess.EventReply, error) {
	log.Println("DispatchEvent called")

	return &coprocess.EventReply{}, nil
}
