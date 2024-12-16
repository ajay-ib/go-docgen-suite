#!/bin/bash

# Set the path to your service
SERVICE_PATH=$(pwd)

# Run the docgen-service to generate documentation
go run -mod=vendor github.com/ajay-ib/go-docgen-suite/cmd/docgen generate --path $SERVICE_PATH --output $SERVICE_PATH/docs