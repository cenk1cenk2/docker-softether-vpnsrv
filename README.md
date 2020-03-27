```
name:         | softether-vpnsrv
compiler:     | docker-compose + dockerfile
version:      | v2.2, 20200327
```

## Description:

A Docker Container that creates a Softether VPN Server instance with "dnsmasq" as DHCP server.

### Features
* S6-Overlay implemented with health checking the server.
* "dnsmasq" as DHCP server.
* Can handle `dnsmasq.conf` with variables.
* Always builds the latest version from the official GitHub repository of SoftEther.
* ~50MB image size, ~15-20MB RAM Usage while standalone.

## Setup:

Clone the GitHub repository to get an environmental variable initiation script and preconfigured docker-compose file if you wish to get a head start. Advised way to run the setup is with docker-compose but it can be run with a long command with docker run.

**Fast Deploy**
* `chmod +x init-env.sh && ./init-env.sh && nano .env` for variables.
* `cp dnsmasq.conf ./cfg/dnsmasq.conf` for "dnsmasq" configuration.
* `cp vpn_server.config ./cfg/vpn_server.config` for vpn server configuration.

`dnsmasq.conf` must include `tap_soft` as tap device both for interface and range, as in the example below.
```
interface=tap_soft
dhcp-range=tap_soft,$SRVIPSUBNET.129,$SRVIPSUBNET.199,255.255.255.0,12h
```

As in the example if $SRVIPSUBNET can be used inside the `dnsmasq.conf` since these instances will be replaced with sed therefore it can be used for various purposes easily changing the configuration.

### Enviroment File
```
# Timezone
TZ=
# VPN Server IP Subnet in form of xx.xx.xx (default: 10.0.0), it can also can rewrite dnsmasq.conf with SED if \$SRVIPSUBNET inside dnsmasq.conf is set."
SRVIPSUBNET=
# Sleep Time for Server Alive Check in Seconds (default: 600)
SLEEPTIME=
# Keep logs or delete them in between sleeptime. To keep set the type to 1.
KEEP_SERVER_LOG=
KEEP_PACKET_LOG=
KEEP_SECURITY_LOG=
```