FROM alpine:latest

ENV SOFTETHER_REPO https://github.com/SoftEtherVPN/SoftEtherVPN.git
ENV SOFTETHER_REPO_API https://api.github.com/repos/SoftEtherVPN/SoftEtherVPN/releases/latest
ENV S6_VERSION 1.21.8.0

# Install s6 overlay
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_VERSION}/s6-overlay-amd64.tar.gz /tmp/
RUN tar xzf /tmp/s6-overlay-amd64.tar.gz -C /

RUN apk --no-cache --no-progress update && \
  apk --no-cache --no-progress upgrade && \
  # Install s6 supervisor
  apk --no-cache --no-progress add bash && \
  mkdir -p /etc/services.d && mkdir /s6-bin && \
  # Install build dependencies
  apk --no-cache --no-progress --virtual .build-deps add git libgcc libstdc++ gcc musl-dev libc-dev g++ ncurses-dev \
  readline-dev openssl-dev cmake make zlib-dev && \
  # Grab and build Softether from GitHub
  git clone ${SOFTETHER_REPO} /tmp/softether && \
  # Get latest Tag
  cd /tmp/softether && export TRACK_LATEST=$(curl --silent ${SOFTETHER_REPO_API} | grep -Po '"tag_name": "\K.*?(?=")') && \
  # Checkout Latest Tag
  git checkout ${TRACK_LATEST} && git submodule init && git submodule update && export USE_MUSL=YES && \
  # Build
  ./configure && make --silent -C tmp && make --silent -C tmp install &&  \
  cp /tmp/softether/build/libcedar.so /tmp/softether/build/libmayaqua.so /usr/lib && \
  # Removing build extensions
  apk del .build-deps && apk del --no-cache --purge && \
  rm -rf /tmp/softether && rm -rf /var/cache/apk/*  && \
  # Deleting unncessary extensions
  rm -rf /usr/local/bin/vpnbridge \
  /usr/local/libexec/softether/vpnbridge && \
  # Reintroduce necassary libraries
  apk add --no-cache --virtual .run-deps \
  libcap libcrypto1.1 libssl1.1 ncurses-libs readline su-exec zlib-dev dhclient && \
  # Link Libraries to Binary
  mkdir /s6-bin/softether-vpnsrv && \
  ln -s /usr/local/bin/vpnserver /s6-bin/softether-vpnsrv/vpnserver && \
  ln -s /usr/local/bin/vpncmd /s6-bin/softether-vpnsrv/vpncmd && \
  ln -s /usr/local/libexec/softether/vpnserver/ /etc && \
  # Install dnsmasq and create the tun network adapter.
  apk add --no-cache dnsmasq iptables && \
  echo "tun" >> /etc/modules && \
  echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf && \
  mkdir -p /dev/net && \
  mknod /dev/net/tun c 10 200

# Open necassary ports
EXPOSE 1443/tcp 992/tcp 1194/tcp 1194/udp 5555/tcp 500/udp 4500/udp 1701/udp

# Move s6 supervisor files inside the container
ADD ./s6 /etc/services.d

ENTRYPOINT ["/init"]