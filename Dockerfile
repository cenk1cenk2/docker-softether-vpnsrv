FROM alpine:latest

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

ENV REPOSITORY https://github.com/SoftEtherVPN/SoftEtherVPN.git

RUN \
  apk --no-cache --no-progress update && \
  apk --no-cache --no-progress upgrade && \
  # Install supervisor
  apk --no-cache --no-progress add tini && \
  mkdir -p /conf && \
  # Install build dependencies
  apk --no-cache --no-progress --virtual .build-deps add git libgcc libstdc++ gcc musl-dev libc-dev g++ ncurses-dev libsodium-dev \
  readline-dev openssl-dev cmake make zlib-dev && \
  # Grab and build Softether from GitHub
  git clone ${REPOSITORY} /tmp/softether && \
  cd /tmp/softether && \
  # init submodules
  git submodule init && git submodule update && \
  # Build
  if [ "${TARGETARCH}" == "arm64" ]; then echo "Forcing to use neon on arm64 platform manually."; sed -ir 's|set(BLAKE2_SRC_PATH $<IF:$<BOOL:${HAS_SSE2}>,${TOP_DIRECTORY}/3rdparty/BLAKE2/sse,${TOP_DIRECTORY}/3rdparty/BLAKE2/ref>)|set(BLAKE2_SRC_PATH ${TOP_DIRECTORY}/3rdparty/BLAKE2/neon)|' ./src/Cedar/CMakeLists.txt; sed -ir 's|set(BLAKE2_SRC $<IF:$<BOOL:${HAS_SSE2}>,${BLAKE2_SRC_PATH}/blake2s.c,${BLAKE2_SRC_PATH}/blake2s-ref.c>)|set(BLAKE2_SRC ${BLAKE2_SRC_PATH}/blake2s-neon.c)|' ./src/Cedar/CMakeLists.txt ; fi && \
  export USE_MUSL=YES && ./configure && make --silent -C build && make --silent -C build install &&  \
  cp /tmp/softether/build/libcedar.so /tmp/softether/build/libmayaqua.so /usr/lib && \
  # Removing build extensions
  apk del .build-deps && apk del --no-cache --purge && \
  rm -rf /tmp/softether && rm -rf /var/cache/apk/*  && \
  # Deleting unncessary extensions
  rm -rf /usr/local/bin/vpnbridge \
  /usr/local/libexec/softether/vpnbridge && \
  # Reintroduce necassary runtime libraries
  apk add --no-cache --virtual .run-deps \
  libcap libcrypto1.1 libssl1.1 ncurses-libs readline su-exec zlib-dev dhclient libsodium-dev && \
  # Link Libraries to Binary
  ln -s /usr/local/bin/vpnserver /usr/bin/softether-vpnsrv && \
  ln -s /usr/local/bin/vpncmd /usr/bin/softether-vpncmd && \
  ln -s /usr/local/libexec/softether/vpnserver/ /etc/softether && \
  # Install dnsmasq and create the tun network adapter.
  apk add --no-cache dnsmasq iptables && \
  echo "tun" >> /etc/modules && \
  echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf && \
  mkdir -p /dev/net && \
  mknod /dev/net/tun c 10 200 && \
  rm -rf /tmp/*

# Advise to open necassary ports
EXPOSE 1443/tcp 992/tcp 1194/tcp 1194/udp 5555/tcp 500/udp 4500/udp 1701/udp

# Adviste to create necassary volumes
VOLUME [ "/conf" ]

# Move host file system
COPY ./.docker/hostfs /
COPY ./dist/pipe-${TARGETOS}-${TARGETARCH}${TARGETVARIANT} /usr/bin/pipe

RUN chmod +x /usr/bin/pipe

# Set working directory back
WORKDIR /

ENTRYPOINT ["tini", "pipe"]
