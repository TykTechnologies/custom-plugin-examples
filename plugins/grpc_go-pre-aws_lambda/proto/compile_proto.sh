#!/bin/bash

echo "Generating bindings for Go."
protoc -I. --go_out=plugins=grpc:./go *.proto
