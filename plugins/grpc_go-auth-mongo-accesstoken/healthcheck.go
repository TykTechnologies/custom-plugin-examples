package main

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthChecker struct{}

func (s *HealthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	logrus.Info("Serving the Check request for health check")
	err := mongoClient.Ping(context.TODO(), nil)
	if err != nil { // timeout case will not trigger this block
		logrus.Fatal(err)
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, err
	}
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	logrus.Info("Serving the Watch request for health check")
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}
