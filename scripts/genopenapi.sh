#!/usr/bin/env bash

set -euo pipefail

shopt -s globstar

if ! [[ "$0" =~ scripts/genopenapi.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

source ./scripts/lib.sh

OPENAPI_ROOT="./api/openapi"

GEN_SERVER=(
#  "chi-server"
#  "echo-server"
  "gin-server"
)

if [ "${#GEN_SERVER[@]}" -ne 1 ];then
  log_error "GEN_SERVER enables more than 1 server,please check."
  exit 225
fi
# ${#GEN_SERVER[@]} 表示 数组长度。
# -ne 1 是 Bash 中的整数比较运算符，表示 “not equal to 1（不等于 1）”。

log_callout "Using ${GEN_SERVER[0]}"



function openapi_files {
  openapi_files=$(ls ${OPENAPI_ROOT})
  echo "${openapi_files[@]}"
}

# output_dir,package_name,service_name
function gen() {
    local output_dir=$1
    local package=$2
    local service=$3

    run mkdir -p "$output_dir"
    run find "$output_dir" -type f -name "*.gen.go" -delete

    prepare_dir "internal/common/client/$service"

    run oapi-codegen -generate types -o "$output_dir/openapi_types.gen.go" -package "$package" "api/openapi/$service.yml"
    run oapi-codegen -generate "$GEN_SERVER" -o "$output_dir/openapi_api.gen.go" -package "$package" "api/openapi/$service.yml"

    run oapi-codegen -generate client -o "internal/common/client/$service/openapi_client.gen.go" -package "$service" "api/openapi/$service.yml"
    run oapi-codegen -generate types -o "internal/common/client/$service/openapi_types.gen.go" -package "$service" "api/openapi/$service.yml"
}
# 这里的 local output_dir=$1 ... 代表第一、第二、第三个参数


gen internal/order/ports ports order

log_success "openapi_generate_success!"