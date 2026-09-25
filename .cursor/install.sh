#!/usr/bin/env bash
# Idempotent Cloud Agent bootstrap for the qrurl monorepo.
# - Installs Go 1.27 (server toolchain) if missing/outdated.
# - Installs the Google Cloud CLI + Cloud Firestore emulator (local dev DB).
# - Downloads Go module and npm dependencies.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GO_VERSION="1.27.0"
GCLOUD_DIR="${HOME}/google-cloud-sdk"

echo "==> Ensuring Go ${GO_VERSION}"
if ! /usr/local/go/bin/go version 2>/dev/null | grep -q "go${GO_VERSION}"; then
  tmp="$(mktemp -d)"
  curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o "${tmp}/go.tgz"
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf "${tmp}/go.tgz"
  rm -rf "${tmp}"
fi
export PATH="/usr/local/go/bin:${PATH}"
go version

echo "==> Ensuring Google Cloud CLI + Firestore emulator"
if [ ! -x "${GCLOUD_DIR}/bin/gcloud" ]; then
  tmp="$(mktemp -d)"
  curl -fsSL "https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-linux-x86_64.tar.gz" -o "${tmp}/gcloud.tgz"
  tar -C "${HOME}" -xzf "${tmp}/gcloud.tgz"
  "${GCLOUD_DIR}/install.sh" --quiet --path-update false --usage-reporting false
  rm -rf "${tmp}"
fi
if ! "${GCLOUD_DIR}/bin/gcloud" components list --only-local-state --format="value(id)" 2>/dev/null | grep -q "cloud-firestore-emulator"; then
  "${GCLOUD_DIR}/bin/gcloud" components install cloud-firestore-emulator beta --quiet
fi

echo "==> Downloading Go modules"
(cd "${REPO_ROOT}/server" && go mod download)

echo "==> Installing client dependencies"
(cd "${REPO_ROOT}/client" && npm ci)

echo "==> Bootstrap complete"
