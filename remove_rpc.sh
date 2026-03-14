#!/bin/bash
set -e

echo "Removing gRPC, NATS, and AMQP from the project..."

# 1. Delete directories
rm -rf docs/proto
rm -rf internal/controller/grpc
rm -rf internal/controller/nats_rpc
rm -rf internal/controller/amqp_rpc
rm -rf pkg/grpcserver
rm -rf pkg/nats
rm -rf pkg/rabbitmq

# 2. Patch internal/app/app.go
sed -i '/\/\/ RabbitMQ RPC Server/,+6d' internal/app/app.go
sed -i '/\/\/ NATS RPC Server/,+6d' internal/app/app.go
sed -i '/\/\/ gRPC Server/,+2d' internal/app/app.go

sed -i '/rmqServer.Start()/d' internal/app/app.go
sed -i '/natsServer.Start()/d' internal/app/app.go
sed -i '/grpcServer.Start()/d' internal/app/app.go

sed -i '/case err = <-grpcServer.Notify():/,+1d' internal/app/app.go
sed -i '/case err = <-rmqServer.Notify():/,+1d' internal/app/app.go
sed -i '/case err = <-natsServer.Notify():/,+1d' internal/app/app.go

sed -i '/err = grpcServer.Shutdown()/,+3d' internal/app/app.go
sed -i '/err = rmqServer.Shutdown()/,+3d' internal/app/app.go
sed -i '/err = natsServer.Shutdown()/,+3d' internal/app/app.go

sed -i '/amqprpc/d' internal/app/app.go
sed -i '/controller\/grpc/d' internal/app/app.go
sed -i '/natsrpc/d' internal/app/app.go
sed -i '/pkg\/grpcserver/d' internal/app/app.go
sed -i '/pkg\/nats\/nats_rpc\/server/d' internal/app/app.go
sed -i '/pkg\/rabbitmq\/rmq_rpc\/server/d' internal/app/app.go

# 3. Patch integration-test/integration_test.go
# Remove test implementations sequentially located at the end of the file FIRST
sed -i '/^\/\/ gRPC Client V1: GetHistory./,$d' integration-test/integration_test.go

sed -i '/docs\/proto\/v1/d' integration-test/integration_test.go
sed -i '/pkg\/nats\/nats_rpc\/client/d' integration-test/integration_test.go
sed -i '/pkg\/rabbitmq\/rmq_rpc\/client/d' integration-test/integration_test.go
sed -i '/google.golang.org\/grpc/d' integration-test/integration_test.go

sed -i '/^[ \t]*\/\/ gRPC/,+1d' integration-test/integration_test.go
sed -i '/^[ \t]*\/\/ RPC configs/,+2d' integration-test/integration_test.go
sed -i '/^[ \t]*\/\/ RabbitMQ RPC/d' integration-test/integration_test.go
sed -i '/rmqURL =/d' integration-test/integration_test.go
sed -i '/natsURL =/d' integration-test/integration_test.go

# Remove unused variables from the tests to fix go lints
sed -i '/^[ \t]*requests.*10/d' integration-test/integration_test.go
sed -i '/^[ \t]*expectedOriginal.*текст/d' integration-test/integration_test.go
sed -i '/^[ \t]*\/\/ Test data/d' integration-test/integration_test.go

# 4. Patch docker-compose.yml
sed -i '/^  rabbitmq:/,/^[ \t]*-[ \t]*rabbitmq\.lvh\.me/d' docker-compose.yml
sed -i '/^  nats:/,/^[ \t]*-[ \t]*nats\.lvh\.me/d' docker-compose.yml

sed -i '/- rabbitmq/d' docker-compose.yml
sed -i '/- nats/d' docker-compose.yml

sed -i '/x-rabbitmq-variables:/,/RABBITMQ_DEFAULT_PASS:/d' docker-compose.yml
sed -i '/# gRPC/,+1d' docker-compose.yml
sed -i '/# NATS/,+2d' docker-compose.yml
sed -i '/# RMQ/,+3d' docker-compose.yml

sed -i '/- "8081:8081"/d' docker-compose.yml
sed -i '/rabbitmq_data:/d' docker-compose.yml
sed -i '/nats_data:/d' docker-compose.yml

# 5. Patch Makefile
sed -i '/proto-v1:/,+6d' Makefile
sed -i 's/ rabbitmq nats//g' Makefile
sed -i 's/ proto-v1//g' Makefile

# 6. Patch nginx/nginx.conf removing upstreams and servers
awk '
BEGIN { skip = 0; buffer = ""; in_server = 0; depth = 0 }
/^[ \t]*server[ \t]*$/ {
    buffer = $0 "\n"
    in_server = 1
    depth = 0
    next
}
in_server {
    buffer = buffer $0 "\n"
    if (/{/) depth++
    if (/}/) depth--
    if (/server_name (grpc|rabbitmq|nats)\.lvh\.me;/) skip = 1
    if (depth == 0 && /}/) {
        if (!skip) printf "%s", buffer
        in_server = 0
        skip = 0
        buffer = ""
    }
    next
}
{ print }
' nginx/nginx.conf >nginx/nginx.conf.tmp && mv nginx/nginx.conf.tmp nginx/nginx.conf

awk '
BEGIN { skip = 0 }
/^[ \t]*upstream (grpc|rabbitmq|nats)[ \t]*$/ { skip = 1; next }
skip {
    if (/^[ \t]*\}[ \t]*$/) skip = 0
    next
}
{ print }
' nginx/nginx.conf >nginx/nginx.conf.tmp && mv nginx/nginx.conf.tmp nginx/nginx.conf

# 7. Patch go.mod and cleanup
sed -i '/google.golang.org\/grpc\/cmd\/protoc-gen-go-grpc/d' go.mod
sed -i '/google.golang.org\/protobuf\/cmd\/protoc-gen-go/d' go.mod
go mod tidy
go fmt ./...
go fix ./...
gofumpt -l -w .
gci write . --skip-generated -s standard -s default

echo "Done! gRPC, NATS, and AMQP and all traces have been cleanly removed."
