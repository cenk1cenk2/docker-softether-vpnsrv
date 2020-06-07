# softether-vpnsrv

[![Build Status](https://drone.kilic.dev/api/badges/cenk1cenk2/softether-vpnsrv/status.svg)](https://drone.kilic.dev/cenk1cenk2/softether-vpnsrv)
![Docker Pulls](https://img.shields.io/docker/pulls/cenk1cenk2/softether-vpnsrv)
![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/cenk1cenk2/softether-vpnsrv)
![Docker Image Version (latest by date)](https://img.shields.io/docker/v/cenk1cenk2/softether-vpnsrv)
![GitHub last commit](https://img.shields.io/github/last-commit/cenk1cenk2/softether-vpnsrv)

```
name:         | softether-vpnsrv
compiler:     | docker-compose + dockerfile
version:      | v5.01.9674, 20200601 | Autoupdated
```

## Description:

SoftEther VPN is free open-source, cross-platform, multi-protocol VPN client and VPN server software, developed as part of Daiyuu Nobori's master's thesis research at the University of Tsukuba. VPN protocols such as SSL VPN, L2TP/IPsec, OpenVPN, and Microsoft Secure Socket Tunneling Protocol are provided in a single VPN server.

This container runs a SoftEther VPN Server bundled together with a DnsMasq DHCP server to distrubute the IPs. With this way it utilizes a Linux virtual ethernet tap device to distrubute the network traffic through.

### Features
#### Always Alive
s6-overlay implemented. SLEEPTIME (default: 3600s) variable can be set through environment variables to check the server, dnsmasq.

#### Graceful shutdown.
Cleans up all the created veth interfaces and undoes all the system changes.

#### Configurable
Can handle `dnsmasq.conf` with variables. SRVIPSUBNET (default: 10.0.0) can be set through environment variables to configure the server at startup.

```bash
### Example and default configuration
port=0
interface=tap_soft
dhcp-option=3
dhcp-option=6
dhcp-range=tap_soft,$SRVIPSUBNET.129,$SRVIPSUBNET.199,255.255.255.0,12h
```

#### Up-to-Date
Dockerfile always pulls the latest tag from the official repository and builds it from scratch.

It will automatically check for updates and tag them matching with the offical repository once every month.

#### Resource-Efficient
Build on top of Alpine linux as base, ~30MB image size, ~15-20MB RAM Usage while standalone.

## Initial Setup
If you dont have Dnsmasq and SoftEther configuration and containerizing existing application. You can use the defaults.

Remember since it creates a veth in the network workspace it has to run in Docker ```--privileged``` mode since it seems that NET_ADMIN capabilities are not enough.

### Dnsmasq Setup
Configuration has defaults as follows. If you mount an empty folder to `/cfg` it will copy these settings over to host file-system.

Default settings:
* Server distrubutes IP addresses from 10.0.0.0/24 subnet.
* IP range is between 10-255.
* Trafic will be tunneled through.

### SoftEther Setup
Configuration has defaults. If you mount an empty folder to `/cfg` it will copy these settings over to host file-system.

* Default port at startup if there is no config file is specified will be 1443.
* Default bridge device is set through the default config file.
* Please check out the normal process for [SoftEther Setup](https://www.softether.org/4-docs/2-howto/9.L2TPIPsec_Setup_Guide_for_SoftEther_VPN_Server/1.Setup_L2TP%2F%2F%2F%2FIPsec_VPN_Server_on_SoftEther_VPN_Server). This can be configured through using the GUI or the CLI.

**Please remember that at initial startup there is no user defined and no admin password for managing server, it is very crucial to set them both ASAP.**

### Command Line Interface
Command line interface can be accessed through `/s6-bin/softether-vpnsrv/vpncmd`.

## Setup:

Clone the GitHub repository to get an environmental variable initiation script and preconfigured docker-compose file if you wish to get a head start. Advised way to run the setup is with docker-compose but it can be run with a long command with docker run.

**Fast Deploy**
```
# Clone repo
git clone git@github.com:cenk1cenk2/softether-vpnsrv.git
# Initiate environment variables for convience
chmod +x init-env.sh
./init-env.sh
nano .env | vi .env

# Create your own configuration or copy existing
cp dnsmasq.conf ./cfg/dnsmasq.conf # Has a default
cp vpn_server.config ./cfg/vpn_server.config # Has a default
```

`dnsmasq.conf` must include `tap_soft` as tap device both for interface and range, as in the example below.
```
interface=tap_soft
dhcp-range=tap_soft,$SRVIPSUBNET.129,$SRVIPSUBNET.199,255.255.255.0,12h
```

### Deploy via Docker
```
docker create \
  --name=softether-vpnsrv \
  -e TZ=Europe/Vienna \
  -e SRVIPSUBNET=10.0.0 \
  -p 1443:1443/tcp \
  -p 992:992/tcp \
  -p 5555:5555/tcp \
  -p 1194:1194/udp \
  -p 500:500/udp \
  -p 4500:4500/udp \
  -p 1701:1701/tcp \
  -v /cfg/vpn_server.config:/cfg/vpn_server.config \
  -v /cfg/dnsmasq.conf:/cfg/dnsmasq.conf \
  --restart unless-stopped \
  --privileged \
  cenk1cenk2/softether-vpnsrv:latest
```

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
