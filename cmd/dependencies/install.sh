#!/bin/bash
set -e

FLAG_FILE=".dependencies_installed"

if [ -f "$FLAG_FILE" ]; then
    echo "Dependencies already installed, skipping."
    exit 0
fi

if [ ! -f "ubuntu_dependencies.txt" ]; then
    echo "File ubuntu_dependencies.txt not found!"
    exit 1
fi

sudo apt update
xargs -a ubuntu_dependencies.txt sudo apt install -y

install_go_tool() {
    local tool_path="$1"   # например, github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    local cmd_name=$(basename "$(echo "$tool_path" | cut -d'@' -f1)")
    echo "Installing $cmd_name via go install $tool_path ..."
    go install "$tool_path"
    if ! command -v "$cmd_name" &> /dev/null; then
        GOPATH=$(go env GOPATH)
        if [ -f "$GOPATH/bin/$cmd_name" ]; then
            sudo ln -sf "$GOPATH/bin/$cmd_name" "/usr/local/bin/$cmd_name"
        else
            echo "Installation of $cmd_name failed"
            exit 1
        fi
    else
        echo "$cmd_name already available"
    fi
}

if [ -f "go_dependencies.txt" ]; then
    echo "Installing Go dependencies from go_dependencies.txt..."
    while IFS= read -r line || [ -n "$line" ]; do
        if [[ -z "$line" || "$line" == \#* ]]; then
            continue
        fi
        install_go_tool "$line"
    done < "go_dependencies.txt"
else
    echo "go_dependencies.txt not found, skipping Go tools installation."
fi

touch "$FLAG_FILE"
echo "All dependencies installed."
