#!/bin/bash

# 检查是否提供了参数
if [ -z "$1" ]; then
  echo "Usage: $0 <file_name_without_extension>"
  exit 1
fi

# 获取传入的参数，作为文件名
FILE=$1

# 使用 protoc 编译指定的 .proto 文件
protoc -I=./proto --go_out=./ --go-grpc_out=./ proto/$FILE.proto
