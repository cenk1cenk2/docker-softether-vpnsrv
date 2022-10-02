# docker-softether-vpnsrv

Initiates the SoftEtherVPN server that will run in this container.

`docker-softether-vpnsrv [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` | &#34;info&#34; |

### Dhcp-Server

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$DHCP_SERVER_TEMPLATE` | Template location for the DHCP server. | `String` | `false` | &#34;/etc/template/dnsmasq.conf.tmpl&#34; |
| `$DHCP_SERVER_LEASE` | DHCP server lease time for clients. | `String` | `false` | &#34;12h&#34; |
| `$DHCP_SERVER_SEND_GATEWAY` | Whether to send the default gateway to the client. Sometimes you do not want to proxy traffic through the network, rather just establish a connection to the VPN network. | `Bool` | `false` | true |
| `$DHCP_SERVER_GATEWAY` | Set the gateway option for the underlying DNS server.  | `String`<br/>dynamic(&#34;cidr address start&#34;) | `false` |  |
| `$DHCP_SERVER_FORWARDING_ZONE` | Set forwarding-zone DNS addresses for the DHCP server. | `StringSlice` | `false` | [8.8.8.8 8.8.4.4] |

### Health

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HEALTH_CHECK_INTERVAL` | Health check interval to the upstream server in duration. | `String` | `false` | &#34;1h&#34; |
| `$HEALTH_DHCP_SERVER_ADDRESS` | Upstream DHCP server address for doing health checks.  | `String`<br/>dynamic(&#34;cidr address start&#34;) | `false` |  |
| `$HEALTH_ENABLE_PING` | Whether to enable the ping check or not. | `Bool` | `false` | true |

### Linux-Bridge

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LINUX_BRIDGE_INTERFACE_NAME` | Interface name for the resulting communication bridge interface. | `String` | `false` | &#34;br100&#34; |
| `$LINUX_BRIDGE_UPSTREAM_INTERFACE` | Interface name for the upstream parent network interface to bridge to, this interface should provide a DHCP server to handle the clients. | `String` | `false` | &#34;eth0&#34; |
| `$LINUX_BRIDGE_USE_DHCP` | Use the upstream DHCP server to get ip for the bridge interface. | `Bool` | `false` | true |
| `$LINUX_BRIDGE_STATIC_IP` | Use a static IP for the bridge interface. | `String` | `false` |  |

### Server

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$SERVER_MODE` | Server mode changes the behavior of the container.  | `String`<br/>enum(&#34;dhcp&#34;, &#34;bridge&#34;) | `true` |  |
| `$SERVER_CIDR_ADDRESS` | CIDR address of the server. | `String` | `false` | &#34;10.0.0.0/24&#34; |

### SoftEther

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$SOFTETHER_TEMPLATE` | Template location for the SoftEtherVPN server. | `String` | `false` | &#34;/etc/template/vpn_server.config.tmpl&#34; |
| `$SOFTETHER_TAP_INTERFACE` | Interface name for SoftEther and the server to bind to as a tap device. | `String` | `false` | &#34;soft&#34; |
| `$SOFTETHER_DEFAULT_HUB` | Default hub name for SoftEtherVPN server. | `String` | `false` | &#34;DEFAULT&#34; |
