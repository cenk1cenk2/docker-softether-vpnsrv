# cenk1cenk2/softether-vpnsrv

[![pipeline status](https://gitlab.kilic.dev/docker/softether-vpnsrv/badges/master/pipeline.svg)](https://gitlab.kilic.dev/docker/softether-vpnsrv/-/commits/master) [![Docker Pulls](https://img.shields.io/docker/pulls/cenk1cenk2/softether-vpnsrv)](https://hub.docker.com/repository/docker/cenk1cenk2/softether-vpnsrv) [![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/cenk1cenk2/softether-vpnsrv)](https://hub.docker.com/repository/docker/cenk1cenk2/softether-vpnsrv) [![Docker Image Version (latest by date)](https://img.shields.io/docker/v/cenk1cenk2/softether-vpnsrv)](https://hub.docker.com/repository/docker/cenk1cenk2/softether-vpnsrv) [![GitHub last commit](https://img.shields.io/github/last-commit/cenk1cenk2/softether-vpnsrv)](https://github.com/cenk1cenk2/softether-vpnsrv)

## Breaking Changes

**This repository contains breaking changes with older versions that are built on top of `s6-overlay`!**

Sorry for any inconvenience caused. It was a decision between making some breaking changes to expand this container or just abandoning this one and creating a new one.

Since my use case now also requires creating bridge networks to access a local network, I had to introduce a new way of configuring this container. This renders the older versions intentionally incompatible for any user to see that something is wrong and hopefully visit here again to easily fix it.

## Description

SoftEther VPN is a free open-source, cross-platform, multi-protocol VPN client and VPN server software developed as part of Daiyuu Nobori's master's thesis research at the University of Tsukuba. VPN protocols such as SSL VPN, L2TP/IPsec, OpenVPN, and Microsoft Secure Socket Tunneling Protocol are provided in a single VPN server.

This container runs a SoftEther VPN Server bundled together with a configuration manager that enables to use of either a DNSMASQ DHCP server to distribute the IPs or bridging to an existing network interface as upstream. It utilizes a Linux virtual ethernet TAP device to distribute the network traffic.

[Read more](https://www.softether.org/) about SoftEther in the official documentation.

<!-- toc -->

- [Features](#features)
  - [Efficient](#efficient)
  - [Up-to-Date](#up-to-date)
  - [Versioning](#versioning)
  - [Health Check](#health-check)
  - [Graceful Shutdown](#graceful-shutdown)
- [Environment Variables](#environment-variables)
  - [General](#general)
  - [Server](#server)
  - [SoftEther](#softether)
  - [Health](#health)
  - [DHCP-Server](#dhcp-server)
  - [Linux-Bridge](#linux-bridge)
- [Setup](#setup)
  - [Volumes](#volumes)
    - [SoftEther Configuration](#softether-configuration)
    - [DNSMASQ Configuration](#dnsmasq-configuration)
    - [Logs](#logs)
  - [Ports](#ports)
  - [Server Mode](#server-mode)
    - [DHCP](#dhcp)
    - [Bridge](#bridge)
  - [Permissions](#permissions)
- [Deploy](#deploy)
- [Interface](#interface)
  - [Command Line Interface](#command-line-interface)
- [SoftEther VPN Client](#softether-vpn-client)

<!-- tocstop -->

---

## Features

### Efficient

Build on top of Alpine Linux as a base, ~30MB image size, ~15-20MB RAM Usage while standby.

### Up-to-Date

This repository is always up to date tracking the [default](https://github.com/SoftEtherVPN/SoftEtherVPN) branch of SoftEtherVPN repository on GitHub. It checks the main repository monthly since there are no frequent updates anymore, and if a new release has been matched it will trigger the build process.

It always builds the application from the source, and while doing that the dependencies will also be updated.

### Versioning

The Docker images are tagged with matching versions to the original repository.

`latest` tag is reserved for the current snapshot of the default branch on the upstream SoftEtherVPN repository.

### Health Check

Periodical health check for monitoring the processes as well as pinging the DHCP server has been implemented. If one of the health checks fails the container will terminate itself. You can use the `docker` restart feature to start the container back on failure.

### Graceful Shutdown

At shutdown or crashes, the container cleans up all the created virtual ethernet interfaces, TAP devices, and undoes all the system changes.

This is a best-effort process and it can not guarantee to finish this process successfully.

## Environment Variables

### General

| Environment Variable | Description                   | Default | Format                                                      | Required |
| -------------------- | ----------------------------- | ------- | ----------------------------------------------------------- | -------- |
| `TZ`                 | Timezone for the server.      |         | string                                                      |          |
| `LOG_LEVEL`          | Log level for the supervisor. | INFO    | enum("PANIC", "FATAL", "WARNING", "INFO", "DEBUG", "TRACE") |          |

### Server

| Environment Variable  | Description                                        | Default       | Format                 | Required |
| --------------------- | -------------------------------------------------- | ------------- | ---------------------- | -------- |
| `SERVER_MODE`         | Server mode changes the behavior of the container. |               | enum("dhcp", "bridge") | yes      |
| `SERVER_CIDR_ADDRESS` | CIDR address of the server.                        | "10.0.0.0/24" | string(cidr)           |          |

### SoftEther

| Environment Variable      | Description                                                             | Default                                | Format       | Required |
| ------------------------- | ----------------------------------------------------------------------- | -------------------------------------- | ------------ | -------- |
| `SOFTETHER_TAP_INTERFACE` | Interface name for SoftEther and the server to bind to as a tap device. | "soft"                                 | string       |          |
| `SOFTETHER_TEMPLATE`      | Template location for the SoftEtherVPN server.                          | "/etc/template/vpn_server.config.tmpl" | string(path) |          |
| `SOFTETHER_DEFAULT_HUB`   | Default hub name for SoftEtherVPN server.                               | "DEFAULT"                              | string       |          |

### Health

| Environment Variable         | Description                                               | Default            | Format           | Required |
| ---------------------------- | --------------------------------------------------------- | ------------------ | ---------------- | -------- |
| `HEALTH_CHECK_INTERVAL`      | Health check interval to the upstream server in duration. | "1h"               | string(duration) |          |
| `HEALTH_DHCP_SERVER_ADDRESS` | Upstream DHCP server address for doing health checks.     | cidr address start | string(cidr)     |          |

### DHCP-Server

The following options are only valid whenever `SERVER_MODE` is `dhcp`.

| Environment Variable | Description | Default | Format | Required |
| --- | --- | --- | --- | --- |
| `DHCP_SERVER_TEMPLATE` | Template location for the DHCP server. | "/etc/template/dnsmasq.conf.tmpl" | string(path) |  |
| `DHCP_SERVER_GATEWAY` | Set the gateway option for the underlying DNS server. | cidr address start | string(cidr) |  |
| `DHCP_SERVER_SEND_GATEWAY` | Whether to send the default gateway to the client. Sometimes you do not want to proxy traffic through the network, rather just establish a connection to the VPN network. |  | boolean |  |
| `DHCP_SERVER_LEASE` | DHCP server lease time for clients. | "12h" | string(duration) |  |
| `DHCP_SERVER_FORWARDING_ZONE` | Set forwarding-zone DNS addresses for the DHCP server. | "8.8.8.8,8.8.4.4" | multiple(string, ",") |  |

### Linux-Bridge

The following options are only valid whenever `SERVER_MODE` is `bridge`.

| Environment Variable | Description | Default | Format | Required |
| --- | --- | --- | --- | --- |
| `LINUX_BRIDGE_INTERFACE_NAME` | Interface name for the resulting communication bridge interface. | "br100" | string |  |
| `LINUX_BRIDGE_UPSTREAM_INTERFACE` | Interface name for the upstream parent network interface to bridge to, this interface should provide a DHCP server to handle the clients. | "eth0" | string |  |

## Setup

**You can mount a persistent folder to the configuration volume at `/conf` to persist configuration data.**

### Volumes

If the configuration for the component is missing, it will auto-generate a new configuration file depending on the environment variables the user has passed.

| Configuration | Location                  |
| ------------- | ------------------------- |
| SoftEtherVPN  | `/conf/vpn_server.config` |
| DNSMASQ       | `/conf/dnsmasq.conf`      |

**It will not manipulate the configuration files if the persistent versions of them are found inside the expected folders already. So if you have a persistent configuration file, it is, unfortunately, your responsibility to get it working.**

#### SoftEther Configuration

The configuration has defaults as follows.

- The default port at startup will be 1443.
- The default bridge device is set through the generation of the configuration file.
- Please check out the normal process for [SoftEther Setup](https://www.softether.org/4-docs/2-howto/9.L2TPIPsec_Setup_Guide_for_SoftEther_VPN_Server/1.Setup_L2TP%2F%2F%2F%2FIPsec_VPN_Server_on_SoftEther_VPN_Server). SoftEtherVPN server can be configured by using the GUI or the CLI.

**Please remember that at the initial startup there is no admin password for managing the server, it is very crucial to set it as soon as possible.**

If you are using a persistent configuration file that is not auto-generated through this container, you should be sure that your HUB configuration is bridged properly with the TAP adapter.

#### DNSMASQ Configuration

The configuration has defaults as follows.

- It will auto-generate a `dnsmasq.conf` depending on the environment variables you have provided.
- If you want to additionally want to add configuration you can always mount a folder to `/etc/dnsmasq.d/` and add your custom configuration.

#### Logs

The log files can be found on `/etc/softether/server_log`, `/etc/softether/security_log`, `/etc/softether/packet_log` inside the container. So you can mount a folder there to obtain the logs from SoftEtherVPN server.

**The auto-generated configuration file will SoftEtherVPN server will have the logging disabled by default.** You can re-enable it through the options of the server.

### Ports

Default ports that SoftEtherVPN server functions on is as follows.

```yaml
- 1443:1443/tcp # softether
- 992:992/tcp # softether alternative
- 5555:5555/tcp # softether alternative
- 1194:1194/udp # openvpn
- 500:500/udp # l2tp IPSec IKE
- 4500:4500/udp # l2tp IPSec
- 1701:1701/tcp # l2tp
```

You can disable the alternatives or the ones that you do not need as you like.

### Server Mode

You can select one of to server modes which are `dhcp` or `bridge`.

#### DHCP

- `dhcp` mode will start a DNSMASQ DHCP server in the background to handle the incoming DHCP requests. This yields much better performance compared to using the `SecureNAT` provided by the SoftEther.
- The configuration file for DNSMASQ will be generated if it is not persistent, with environment variables from the DHCP Server section.
- It will be attached to the given `SOFTETHER_TAP_INTERFACE`. So whenever you do attach a persistent configuration file manually, please be sure that it points to the correct tap interface. The configuration utility will also assign a static IP address to the given interface making it compatible with how DNSMASQ actually starts up.

#### Bridge

- `bridge` mode will create a new bridge adapter and add the TAP adapter and the upstream adapter to it.
- The upstream DHCP server on the upstream adapter interface will be providing the DHCP setup. This will allow you to connect to a local network.
- Performance is abysmal if you are using the bridge adapter or the upstream adapter directly. SoftEtherVPN server behaves the best whenever it goes through a TAP adapter, so that is why there is still a TAP adapter in between.
- Please be sure to set the server CIDR address properly, because this would set the ping health check to the upstream DHCP server which in the case it would not respond will fail the health check of the container.

### Permissions

**This container needs extra Linux capabilities as provided by [thjderjktyrjkt](https://github.com/thjderjktyrjkt), you can find this in the related issue [#20](https://github.com/cenk1cenk2/docker-softether-vpnsrv/issues/20).**

Basically, it needs the following capabilities to function properly, while it creates a virtual network adapter for the communication.

```yaml
cap_add:
  - SETGID
  - SETUID
  - NET_ADMIN
  - NET_RAW
  - NET_BIND_SERVICE
```

It also needs to access the `tun` device of the machine which can be added as follows.

```yaml
devices:
  - /dev/net/tun
```

## Deploy

You can check out the example setup inside `docker-compose.yml` to see how this container can become operational. Please mind the environment variables for the configuration as well as the section about the configuration files and their generation.

## Interface

If you ever want to interact with the underlying applications you can execute the tasks in the container.

### Command Line Interface

Command-line interface in the container can be accessed through `softether-vpncmd`.

## SoftEther VPN Client

The complementary Docker container for this server can be found [here](https://github.com/cenk1cenk2/softether-vpncli).
