package pipe

import (
	"github.com/urfave/cli/v2"
)

const (
	category_health       = "health"
	category_dhcp_server  = "dhcp-server"
	category_linux_bridge = "linux-bridge"
	category_server       = "server"
)

var Flags = []cli.Flag{
	// runtime

	&cli.StringFlag{
		Name:        "health.check-interval",
		Usage:       "Health check interval to the upstream server in duration.",
		Category:    category_health,
		Required:    false,
		EnvVars:     []string{"HEALTH_CHECK_INTERVAL"},
		Value:       "1h",
		Destination: &TL.Pipe.Health.CheckInterval,
	},

	&cli.StringFlag{
		Name:        "health.dhcp-server-address",
		Usage:       "Upstream DHCP server address for doing health checks. [default: cidr address start]",
		Category:    category_health,
		Required:    false,
		EnvVars:     []string{"HEALTH_DHCP_SERVER_ADDRESS"},
		Value:       "",
		Destination: &TL.Pipe.Health.DhcpServerAddress,
	},

	// dhcp server

	&cli.StringFlag{
		Name:        "dhcp-server.template",
		Usage:       "Template location for the DHCP server.",
		Category:    category_dhcp_server,
		Required:    false,
		EnvVars:     []string{"DHCP_SERVER_TEMPLATE"},
		Value:       "/template/dnsmasq.conf.tmpl",
		Destination: &TL.Pipe.DhcpServer.Template,
	},

	&cli.StringFlag{
		Name:        "dhcp-server.lease",
		Usage:       "DHCP server lease time for clients.",
		Category:    category_dhcp_server,
		Required:    false,
		EnvVars:     []string{"DHCP_SERVER_LEASE"},
		Value:       "12h",
		Destination: &TL.Pipe.DhcpServer.Lease,
	},

	&cli.StringFlag{
		Name:        "dhcp-server.tap-interface",
		Usage:       "Interface name for SoftEther and the DNS server to bind to as a tap device.",
		Category:    category_dhcp_server,
		Required:    false,
		EnvVars:     []string{"DHCP_SERVER_TAP_INTERFACE"},
		Value:       "soft",
		Destination: &TL.Pipe.DhcpServer.TapInterface,
	},

	&cli.BoolFlag{
		Name:        "dhcp-server.send-gateway",
		Usage:       "Whether to send the default gateway to the client. Sometimes you do not want to proxy traffic through the network, rather just establish a connection to the VPN network.",
		Category:    category_dhcp_server,
		Required:    false,
		EnvVars:     []string{"DHCP_SERVER_SEND_GATEWAY"},
		Value:       true,
		Destination: &TL.Pipe.DhcpServer.SendGateway,
	},

	&cli.StringFlag{
		Name:        "dhcp-server.gateway",
		Usage:       "Set the gateway option for the underlying DNS server. [default: cidr address start]",
		Category:    category_dhcp_server,
		Required:    false,
		EnvVars:     []string{"DHCP_SERVER_GATEWAY"},
		Value:       "",
		Destination: &TL.Pipe.DhcpServer.Gateway,
	},

	&cli.StringSliceFlag{
		Name:        "dhcp-server.forwarding-zone",
		Usage:       "Set forwarding-zone DNS addresses for the DHCP server.",
		Category:    category_dhcp_server,
		Required:    false,
		EnvVars:     []string{"DHCP_SERVER_FORWARDING_ZONE"},
		Value:       cli.NewStringSlice("8.8.8.8", "8.8.4.4"),
		Destination: &TL.Pipe.DhcpServer.ForwardingZone,
	},

	&cli.StringFlag{
		Name:        "dhcp-server.options",
		Usage:       `Set custom options for the DHCP server. (format: json { "key": "value" })`,
		Category:    category_dhcp_server,
		Required:    false,
		EnvVars:     []string{"DHCP_SERVER_OPTIONS"},
		Value:       "{}",
		Destination: &TL.Pipe.DhcpServer.Options,
	},

	// linux bridge

	&cli.StringFlag{
		Name:        "linux-bridge.bridge-interface",
		Usage:       "Interface name for the resulting communication bridge interface.",
		Category:    category_linux_bridge,
		Required:    false,
		EnvVars:     []string{"LINUX_BRIDGE_INTERFACE"},
		Value:       "br100",
		Destination: &TL.Pipe.LinuxBridge.BridgeInterface,
	},

	&cli.StringFlag{
		Name:        "linux-bridge.upstream-interface",
		Usage:       "Interface name for the upstream parent network interface to bridge to, this interface should provide a DHCP server to handle the clients.",
		Category:    category_linux_bridge,
		Required:    false,
		EnvVars:     []string{"LINUX_UPSTREAM_INTERFACE"},
		Value:       "eth0",
		Destination: &TL.Pipe.LinuxBridge.UpstreamInterface,
	},

	&cli.StringFlag{
		Name:        "linux-bridge.tap-interface",
		Usage:       "Interface name for SoftEther to bind to as a tap device. Directly binding to the upstream interface causes slow connection, so it is a middle ground between the upstream interface and the SoftEtherVPN server.",
		Category:    category_linux_bridge,
		Required:    false,
		EnvVars:     []string{"LINUX_PARENT_INTERFACE"},
		Value:       "eth0",
		Destination: &TL.Pipe.LinuxBridge.UpstreamInterface,
	},

	// server

	&cli.StringFlag{
		Name:        "server.mode",
		Usage:       `Server mode changes the behavior of the container. [enum: "dhcp", "bridge"]`,
		Category:    category_server,
		Required:    true,
		EnvVars:     []string{"SERVER_MODE"},
		Destination: &TL.Pipe.Server.Mode,
	},

	&cli.StringFlag{
		Name:        "server.cidr-address",
		Usage:       `Server mode changes the behavior of the container. [enum: "dhcp", "bridge"]`,
		Category:    category_server,
		Required:    false,
		Value:       "10.0.0.0/24",
		EnvVars:     []string{"SERVER_CIDR_ADDRESS"},
		Destination: &TL.Pipe.Server.CidrAddress,
	},
}
