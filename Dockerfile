FROM golang:1.20-bullseye as builder

ADD . /webtunnel

ENV CGO_ENABLED=0

WORKDIR /webtunnel

RUN go build -ldflags="-s -w" -o "build/server" gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/webtunnel/main/server

FROM debian:bullseye-slim

COPY --from=builder /webtunnel/build/server /usr/bin/webtunnel-server

# Install dependencies to add Tor's repository.
RUN apt-get update && apt-get install -y \
    curl \
    gpg \
    gpg-agent \
    ca-certificates \
    libcap2-bin \
    --no-install-recommends

# See: <https://2019.www.torproject.org/docs/debian.html.en>
RUN curl https://deb.torproject.org/torproject.org/A3C4F0F979CAA22CDBA8F512EE8CBC9E886DDD89.asc | gpg --import
RUN gpg --export A3C4F0F979CAA22CDBA8F512EE8CBC9E886DDD89 | apt-key add -

RUN printf "deb https://deb.torproject.org/torproject.org bullseye main\n" >> /etc/apt/sources.list.d/tor.list

# Install remaining dependencies.
RUN apt-get update && apt-get install -y \
    tor \
    tor-geoipdb \
    --no-install-recommends

# Our torrc is generated at run-time by the script start-tor.sh.
RUN rm /etc/tor/torrc
RUN chown debian-tor:debian-tor /etc/tor
RUN chown debian-tor:debian-tor /var/log/tor

ADD release/container/start-tor.sh /usr/local/bin
RUN chmod 0755 /usr/local/bin/start-tor.sh

ADD release/container/get-bridge-line.sh /usr/local/bin
RUN chmod 0755 /usr/local/bin/get-bridge-line.sh

ENTRYPOINT ["/usr/local/bin/start-tor.sh"]
