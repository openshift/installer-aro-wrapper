#!/usr/bin/env bash
set -euo pipefail
# Download the logging images. This script is executed as a oneshot
# service by systemd, because we cannot make use of Requires and a
# simple service: https://github.com/systemd/systemd/issues/1312.
#
# This script continues trying to download the logging images until
# successful because we cannot use Restart=on-failure with a oneshot
# service: https://github.com/systemd/systemd/issues/2582.
#

echo "Pulling $MDSDIMAGE..."
while true
do
    if podman pull --quiet "$MDSDIMAGE"
    then
        break
    else
        echo "Pull failed. Retrying $MDSDIMAGE..."
    fi
done

echo "Pulling $FLUENTIMAGE..."
while true
do
    if podman pull --quiet "$FLUENTIMAGE"
    then
        break
    else
        echo "Pull failed. Retrying $FLUENTIMAGE..."
    fi
done
