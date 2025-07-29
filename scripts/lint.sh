#!/usr/bin/env bash

set -euo pipefail

source ./scripts/lib.sh

# 安装指定工具（如果未安装）
function install_if_not_exist() {
  TOOL_NAME=$1
  INSTALL_URL=$2
  if command -v "$TOOL_NAME" &> /dev/null; then
    log_callout "$TOOL_NAME is already installed."
  else
    log_cmd "$TOOL_NAME is not installed. Installing..."
    run go install "$INSTALL_URL"
  fi
}

# 1. 安装 go-cleanarch（架构检查工具）
install_if_not_exist go-cleanarch github.com/roblaszczak/go-cleanarch@latest

# 2. 安装或升级 golangci-lint（兼容 Go 最新版本）
log_cmd "Checking golangci-lint version..."
GO_VERSION=$(go version | awk '{print $3}')
log_callout "Current Go version: $GO_VERSION"

if command -v golangci-lint >/dev/null 2>&1; then
  CURRENT_LINT_VERSION=$(golangci-lint --version | awk '{print $4}')
  log_callout "golangci-lint is installed (version: $CURRENT_LINT_VERSION)"
else
  log_cmd "golangci-lint is not installed. Installing latest version..."
fi

# 总是使用最新版，避免版本不兼容 Go
run go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 3. 执行架构检查
run go-cleanarch

# 4. 打印模块名（项目定义）
log_info "lint modules:"
log_info "$(modules)"

# 5. 格式化导入
run goimports -w -l .

# 6. 遍历模块进行 golangci-lint 检查
while read -r module; do
  run cd ./internal/"$module"
  run golangci-lint run --config "$ROOT_DIR/.golangci.yaml"
  run cd -
done < <(modules)
