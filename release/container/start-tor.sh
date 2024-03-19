#!/usr/bin/env bash

NICK=${NICKNAME:-WebTunnelBr}

echo "Using NICKNAME=${NICK}, OR_PORT=${OR_PORT}, PT_PORT=${PT_PORT}, and EMAIL=${EMAIL}."

ADDITIONAL_VARIABLES_PREFIX="WEBTUNNELV_"
ADDITIONAL_VARIABLES=

if [[ "$WEBTUNNEL_ENABLE_ADDITIONAL_VARIABLES" == "1" ]]
then
    ADDITIONAL_VARIABLES="# Additional properties from processed '$ADDITIONAL_VARIABLES_PREFIX' environment variables"
    echo "Additional properties from '$ADDITIONAL_VARIABLES_PREFIX' environment variables processing enabled"

    IFS=$'\n'
    for V in $(env | grep "^$ADDITIONAL_VARIABLES_PREFIX"); do
        VKEY_ORG="$(echo $V | cut -d '=' -f1)"
        VKEY="${VKEY_ORG#$ADDITIONAL_VARIABLES_PREFIX}"
        VVALUE="$(echo $V | cut -d '=' -f2)"
        echo "Overriding '$VKEY' with value '$VVALUE'"
        ADDITIONAL_VARIABLES="$ADDITIONAL_VARIABLES"$'\n'"$VKEY $VVALUE"
    done
fi

cat > /etc/tor/torrc << EOF
RunAsDaemon 0
# We don't need an open SOCKS port.
SocksPort 0
BridgeRelay 1
Nickname ${NICK}
Log notice file /var/log/tor/log
Log notice stdout
ServerTransportPlugin webtunnel exec /usr/bin/webtunnel-server
ExtORPort auto
DataDirectory /var/lib/tor

# The variable "OR_PORT" is replaced with the OR port.
ORPort ${OR_PORT}

# The variable "PT_PORT" is replaced with the obfs4 port.
ServerTransportListenAddr webtunnel 0.0.0.0:${PT_PORT}

# The variable "EMAIL" is replaced with the operator's email address.
ContactInfo ${EMAIL}

ServerTransportOptions webtunnel url=${WEBTUNNEL_URL}

$ADDITIONAL_VARIABLES
EOF

echo "Starting tor."
exec tor -f /etc/tor/torrc
