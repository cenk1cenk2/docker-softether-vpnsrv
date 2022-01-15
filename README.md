# cenk1cenk2/softether-vpnsrv

[![Build Status](https://drone.kilic.dev/api/badges/cenk1cenk2/docker-softether-vpnsrv/status.svg)](https://drone.kilic.dev/cenk1cenk2/docker-softether-vpnsrv) [![Docker Pulls](https://img.shields.io/docker/pulls/cenk1cenk2/softether-vpnsrv)](https://hub.docker.com/repository/docker/cenk1cenk2/softether-vpnsrv) [![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/cenk1cenk2/softether-vpnsrv)](https://hub.docker.com/repository/docker/cenk1cenk2/softether-vpnsrv) [![Docker Image Version (latest by date)](https://img.shields.io/docker/v/cenk1cenk2/softether-vpnsrv)](https://hub.docker.com/repository/docker/cenk1cenk2/softether-vpnsrv) [![GitHub last commit](https://img.shields.io/github/last-commit/cenk1cenk2/softether-vpnsrv)](https://github.com/cenk1cenk2/softether-vpnsrv)

<!-- toc -->

- [Description](#description)
- [Features](#features)
  - [Resource-Efficient](#resource-efficient)
  - [Up-to-Date](#up-to-date)
  - [Version Tracking](#version-tracking)
  - [Always Alive](#always-alive)
  - [Graceful Shutdown](#graceful-shutdown)
- [Environment Variables](#environment-variables)
- [Setup](#setup)
  - [DNSMASQ Setup](#dnsmasq-setup)
    - [Interpolating Variables](#interpolating-variables)
    - [Mounting Custom DNSMASQ Configuration File](#mounting-custom-dnsmasq-configuration-file)
  - [SoftEther Setup](#softether-setup)
    - [Mounting Custom SoftEther Configuration File](#mounting-custom-softether-configuration-file)
- [Deploy](#deploy)
  - [docker-compose](#docker-compose)
  - [docker](#docker)
- [Interface](#interface)
  - [Command Line Interface](#command-line-interface)
- [SoftEther VPN Client](#softether-vpn-client)

<!-- tocstop -->

---

## Description

SoftEther VPN is a free open-source, cross-platform, multi-protocol VPN client and VPN server software developed as part of Daiyuu Nobori's master's thesis research at the University of Tsukuba. VPN protocols such as SSL VPN, L2TP/IPsec, OpenVPN, and Microsoft Secure Socket Tunneling Protocol are provided in a single VPN server.

This container runs a SoftEther VPN Server bundled together with a DNSMASQ DHCP server to distribute the IPs. In this way, it utilizes a Linux virtual ethernet tap device to distribute the network traffic.

[Read more](https://www.softether.org/) about SoftEther in the official documentation.

## Features

### Resource-Efficient

Build on top of Alpine Linux as a base, ~30MB image size, ~15-20MB RAM Usage while standby.

### Up-to-Date

This repository is always up to date tracking the [default](https://github.com/SoftEtherVPN/SoftEtherVPN) repository of SoftEther VPN on GitHub. It checks the main repository monthly since there are no frequent updates anymore, and if a new release has been matched it will trigger the build process.

It always builds the application from the source, and while doing that the dependencies will also be updated.

### Version Tracking

The Docker images are given matching versions to the original repository. If an update has been made on this repository itself, it will append a suffix to the original version.

Please use the `edge` version since tags are few and in-between the build from master goes to the `edge` version.

### Always Alive

[s6-overlay](https://github.com/just-containers/s6-overlay) is implemented to check whether everything is working as expected and do a sanity check with pinging the main VPN server periodically.

An environment variable, namely `SLEEPTIME` can be set in seconds to determine the period of this check.

If the periodic check fails, it will go into graceful shutdown mode and clear any residue like tap devices, virtual network adapters, and such, so it can restart from scratch.

### Graceful Shutdown

At shutdown or crashes, the container cleans up all the created VETH interfaces, tap devices, and undoes all the system changes.

## Environment Variables

| Environment Variable | Description                                                                       | Default Value |
| -------------------- | --------------------------------------------------------------------------------- | ------------- |
| `TZ`                 | Timezone for the server.                                                          |               |
| `LOG_LEVEL`          | Log level for the scripts. Can be: [ SILENT, ERROR, WARN, LIFETIME, INFO, DEBUG ] | INFO          |
| `SLEEPTIME`          | The time in seconds between checks of whether everything is working.              | 3600          |
| `KEEP_SERVER_LOG`    | Keep server logs, set to 1 to keep.                                               |               |
| `KEEP_PACKET_LOG`    | Keep packet logs, set to 1 to keep.                                               |               |
| `KEEP_SECURITY_LOG`  | Keep security logs, set to 1 to keep.                                             |               |
| `SRVIPSUBNET`        | Subnet of the distributed IP addresses by DNSMASQ.                                | 10.0.0        |
| `SRVIPNETMASK`       | Netmask for the subnet.                                                           | 255.255.255.0 |
| `DHCP_START`         | Start address of distributed IP addresses.                                        | 10            |
| `DHCP_END`           | End address of distributed IP addresses.                                          | 254           |
| `DHCP_LEASE`         | Lease time of distributed IP addresses.                                           | 12h           |

## Setup

If you do not have any default configuration for your the defaults will be applied, and the configuration will reside in the `/config` folder.

**You can mount a persistent folder to this folder to further edit and persist both SoftEther and DNSMASQ data. You can also let it generate the defaults by mounting an empty folder to there.**

**Remember since it creates a virtual ethernet in the network workspace it has to run in Docker `--privileged` mode since it seems that NET_ADMIN capabilities are not enough.**

### DNSMASQ Setup

The configuration has defaults as follows.

- Server distributes IP addresses from the 10.0.0.0/24 subnet.
- IP range is between 10-255.
- Traffic will be tunneled through.

#### Interpolating Variables

Can handle `dnsmasq.conf` with variables. If you mount a custom `dnsmasq.conf`, you can still use these variables. It will be interpolated by `sed` while linking it to the actual location.

```bash
### Example and default configuration
port=0
interface=tap_soft
dhcp-range=tap_soft,$SRVIPSUBNET.$DHCP_START,$SRVIPSUBNET.$DHCP_END,$SRVIPNETMASK,$DHCP_LEASE
dhcp-option=tap_soft,3,$SRVIPSUBNET.1
dhcp-option=tap_soft,6,8.8.8.8,8.8.4.4
```

#### Mounting Custom DNSMASQ Configuration File

**For further customization ensure that you have a `dnsmasq.conf` file mounted in `/config` folder.** Your custom `dnsmasq.conf` must include `tap_soft` as tap device both for interface and range, as in the example below.

```conf
interface=tap_soft
dhcp-range=tap_soft,${REST_OF_THE_VARIABLES}
```

### SoftEther Setup

The configuration has defaults as follows.

- Default port at startup will be 1443.
- Default bridge device is set through the default config file.
- Please check out the normal process for [SoftEther Setup](https://www.softether.org/4-docs/2-howto/9.L2TPIPsec_Setup_Guide_for_SoftEther_VPN_Server/1.Setup_L2TP%2F%2F%2F%2FIPsec_VPN_Server_on_SoftEther_VPN_Server). This can be configured through using the GUI or the CLI.

**Please remember that at initial startup there is no user-defined and no admin password for managing the server, it is very crucial to set them both ASAP.**

#### Mounting Custom SoftEther Configuration File

**For further customization ensure that you have a `vpn_server.config` file mounted in the `/config` folder.**

## Deploy

### docker-compose

Clone the GitHub repository to get an environmental variable initiation script and preconfigured docker-compose file if you wish to get a head start. Advised way to run the setup is with docker-compose but it can be run with a long command with docker run.

```bash
# Clone repository
git clone git@github.com:cenk1cenk2/softether-vpnsrv.git

# Initiate environment variables for convenience
chmod +x init-env.sh

./init-env.sh

nvim .env

# Create your own configuration or copy existing
cp dnsmasq.conf ./volumes/softether-vpnsrv/dnsmasq.conf # Has a default
cp vpn_server.config ./volumes/softether-vpnsrv/vpn_server.config # Has a default
```

### docker

```shell
docker create \
  --name=softether-vpnsrv \
  -e TZ=Europe/Vienna \
  -e SRVIPSUBNET=10.0.0 \
  -e SRVIPNETMASK=255.255.255.0 \
  -p 1443:1443/tcp \
  -p 992:992/tcp \
  -p 5555:5555/tcp \
  -p 1194:1194/udp \
  -p 500:500/udp \
  -p 4500:4500/udp \
  -p 1701:1701/tcp \
  -v $PWD/cfg/vpn_server.config:/cfg/vpn_server.config \
  -v $PWD/cfg/dnsmasq.conf:/cfg/dnsmasq.conf \
  --restart unless-stopped \
  --privileged \
  cenk1cenk2/softether-vpnsrv:latest
```

## Interface

If you ever want to interact with the underlying applications you can execute the tasks in the container.

### Command Line Interface

Command line interface can be accessed through `/s6-bin/softether-vpnsrv/vpncmd`.

## SoftEther VPN Client

The complementary Docker container for this server can be found at [here](https://github.com/cenk1cenk2/softether-vpncli).
