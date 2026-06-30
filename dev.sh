#!/bin/bash
set -e

# VisuTask 开发调试脚本
# 启动 Wails v3 dev 模式，文件变更自动重建并重启应用

# ─── 环境配置 ───
GO_VERSION="1.26.2"
GO_SDK="$HOME/sdk/go${GO_VERSION}/bin"
GO_BIN="$HOME/go/bin"
export PATH="$GO_SDK:$GO_BIN:$PATH"

# ─── 环境检查 ───
if ! command -v go &>/dev/null; then
    echo "❌ Go not found at $GO_SDK"
    exit 1
fi

if ! command -v pnpm &>/dev/null; then
    echo "❌ pnpm not found"
    exit 1
fi

if ! command -v wails3 &>/dev/null; then
    echo "Installing Wails v3 CLI..."
    go install github.com/wailsapp/wails/v3/cmd/wails3@latest
fi

# ─── 切换到项目根目录 ───
cd "$(dirname "$0")"

echo "========================================="
echo "  VisuTask Dev Mode (Wails v3)"
echo "  Go:    $(go version)"
echo "  Wails: $(wails3 version 2>&1 | head -1)"
echo "  pnpm:  $(pnpm --version)"
echo "  Dir:   $(pwd)"
echo "========================================="
echo ""
echo "  Auto rebuild on file change:"
echo "    *.go   -> pnpm build -> go build -> restart"
echo "    *.tsx  -> pnpm build -> go build -> restart"
echo "    *.ts   -> pnpm build -> go build -> restart"
echo "    *.css  -> pnpm build -> go build -> restart"
echo ""
echo "  Press Ctrl+C to stop"
echo "========================================="
echo ""

# ─── 启动 Wails dev 模式 ───
wails3 dev -config ./build/config.yml
