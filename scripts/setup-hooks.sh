#!/usr/bin/env bash
# Configure git to use the project's hooks directory.
set -euo pipefail

git config core.hooksPath .githooks
echo "Git hooks configured: .githooks"
