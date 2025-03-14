# WebTunnel

Pluggable Transport based on HTTP Upgrade(HTTPT)

WebTunnel is pluggable transport that attempt to imitate web browsing activities based on [HTTPT](https://censorbib.nymity.ch/#Frolov2020b).

## Client Usage
Connect to a WebTunnel server with a Tor configuration file like:
```
UseBridges 1
DataDirectory datadir

ClientTransportPlugin webtunnel exec ./client

Bridge webtunnel 192.0.2.3:1 url=https://akbwadp9lc5fyyz0cj4d76z643pxgbfh6oyc-167-71-71-157.sslip.io/5m9yq0j4ghkz0fz7qmuw58cvbjon0ebnrsp0

SocksPort auto

Log info
```
## Running a WebTunnel Bridge

You can help censored users connect to the Tor network by running a WebTunnel bridge, see our [community documentation](https://community.torproject.org/relay/setup/webtunnel/) for more details.

## WebTunnel Client Bridgeline Format

#### url: string

`url` determines the HTTP layer host and path.

It should be an HTTPS protocol URL string that points to the server endpoint where the WebTunnel Bridge is hosted.

#### version: string

`version` determines the version of the server. This allows the client to adjust its protocol based on the options supported on the server side.

#### addr: string

`addr` determines the Network Layer (TCP) endpoint of the server. By default, it is the same as the host with the port from the URL. (lyrebird version)

#### servername: string

`servername` determines the Transport Layer Security (TLS) server name indication. By default, it is the same as the host without the port from the URL.

#### utls: enum

`utls` determines the utls tls client hello fingerpint.

valid vlues are:
- `none` : use go's default tls fingerprint
