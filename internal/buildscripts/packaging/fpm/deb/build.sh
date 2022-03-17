#!/bin/bash

set -euxo pipefail

SCRIPT_DIR="$( cd "$( dirname ${BASH_SOURCE[0]} )" && pwd )"
REPO_DIR="$( cd "$SCRIPT_DIR/../../../../../" && pwd )"
VERSION="${1:-}"
ARCH="${2:-"amd64"}"
OUTPUT_DIR="${3:-"$REPO_DIR/dist/"}"
OTELCONTRIBCOL_PATH="$REPO_DIR/bin/otelcontribcol_linux_$ARCH"

. $SCRIPT_DIR/../common.sh

if [[ -z "$VERSION" ]]; then
    latest_tag="$( git describe --abbrev=0 --match v[0-9]* )"
    VERSION="${latest_tag}-post"
fi

mkdir -p "$OUTPUT_DIR"

fpm -s dir -t deb -n $PKG_NAME -v ${VERSION#v} -f -p "$OUTPUT_DIR" \
    --vendor "$PKG_VENDOR" \
    --maintainer "$PKG_MAINTAINER" \
    --description "$PKG_DESCRIPTION" \
    --license "$PKG_LICENSE" \
    --url "$PKG_URL" \
    --architecture "$ARCH" \
    --deb-dist "stable" \
    --deb-user "$PKG_USER" \
    --deb-group "$PKG_GROUP" \
    --before-install "$PREINSTALL_PATH" \
    --after-install "$POSTINSTALL_PATH" \
    --pre-uninstall "$PREUNINSTALL_PATH" \
    $SERVICE_PATH=/lib/systemd/system/$SERVICE_NAME.service \
    $OTELCONTRIBCOL_PATH=/usr/bin/otelcontribcol \
    $CONFIG_PATH=/etc/otel-contrib-collector/config.yaml \
    $ENVFILE_PATH=/etc/otel-contrib-collector/otel-contrib-collector.conf