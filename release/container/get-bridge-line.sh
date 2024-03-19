#!/usr/bin/env bash
#
# This script extracts the pieces that we need to compile our bridge line.
# This will have to do until the following bug is fixed:
# <https://gitlab.torproject.org/tpo/core/tor/-/issues/29128>

TOR_LOG=/var/log/tor/log

if [ ! -r "$TOR_LOG" ]
then
    echo "Cannot read Tor's log file ${TOR_LOG}. This is a bug."
    exit 1
fi

fingerprint=$(grep "Your Tor server's identity key *fingerprint is" "$TOR_LOG" | \
    sed "s/.*\([0-9A-F]\{40\}\)'$/\1/" | \
    tail -1)

imaginaryaddr=$(grep 'Registered server transport' "$TOR_LOG" | sed -E "s/.*?(\[[0-9a-f:]*\]:443)'$/\\1/gm" | \
    tail -1)

echo "webtunnel ${imaginaryaddr} ${fingerprint} url=${WEBTUNNEL_URL}"
