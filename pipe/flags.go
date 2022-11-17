package pipe

import (
	"time"

	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

//revive:disable:line-length-limit

const (
	category_health       = "Health"
	category_dhcp_server  = "DHCP Server"
	category_linux_bridge = "Linux Bridge"
	category_server       = "Server"
	category_softether    = "SoftEther"
)

var Flags = []cli.Flag{
	// runtime

	&cli.DurationFlag{
		Name:        "health.check-interval",
		Usage:       "Health check interval to the upstream server in duration.",
		Category:    category_health,
		Required:    false,
		EnvVars:     []string{"HEALTH_CHECK_INTERVAL"},
		Value:       time.Minute * 10,
		Destination: &TL.Pipe.Health.CheckInterval,
	},

	&cli.StringFlag{
		Name:        "health.dhcp-server-address",
		Usage:       "Upstream DHCP server address for doing health checks.",
		Category:    category_health,
		Required:    false,
		EnvVars:     []string{"HEALTH_DHCP_SERVER_ADDRESS"},
		DefaultText: "CIDR address range start",
		Value:       "",
		Destination: &TL.Pipe.Health.DhcpServerAddress,
	},

	&cli.BoolFlag{
		Name:        "health.enable-ping",
		Usage:       "Whether to enable the ping check or not.",
		Category:    category_health,
		Required:    false,
		EnvVars:     []string{"HEALTH_ENABLE_PING"},
		Value:       true,
		Destination: &TL.Pipe.Health.EnablePing,
	},

	// dhcp server

	&cli.StringFlag{
		Name:        "dhcp-server.template",
		Usage:       "Template location for the DHCP server.",
		Category:    category_dhcp_server,
		Required:    false,
		EnvVars:     []string{"DHCP_SERVER_TEMPLATE"},
		Value:       "/etc/template/dnsmasq.conf.tmpl",
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
		Usage:       "Set the gateway option for the underlying DNS server.",
		Category:    category_dhcp_server,
		Required:    false,
		EnvVars:     []string{"DHCP_SERVER_GATEWAY"},
		DefaultText: "CIDR address range start",
		Value:       "",
		Destination: &TL.Pipe.DhcpServer.Gateway,
	},

	&cli.StringSliceFlag{
		Name:     "dhcp-server.forwarding-zone",
		Usage:    "Set forwarding-zone DNS addresses for the DHCP server.",
		Category: category_dhcp_server,
		Required: false,
		EnvVars:  []string{"DHCP_SERVER_FORWARDING_ZONE"},
		Value:    cli.NewStringSlice("8.8.8.8", "8.8.4.4"),
	},

	// linux bridge

	&cli.StringFlag{
		Name:        "linux-bridge.bridge-interface",
		Usage:       "Interface name for the resulting communication bridge interface.",
		Category:    category_linux_bridge,
		Required:    false,
		EnvVars:     []string{"LINUX_BRIDGE_INTERFACE_NAME"},
		Value:       "br100",
		Destination: &TL.Pipe.LinuxBridge.BridgeInterface,
	},

	&cli.StringFlag{
		Name:        "linux-bridge.upstream-interface",
		Usage:       "Interface name for the upstream parent network interface to bridge to, this interface should provide a DHCP server to handle the clients.",
		Category:    category_linux_bridge,
		Required:    false,
		EnvVars:     []string{"LINUX_BRIDGE_UPSTREAM_INTERFACE"},
		Value:       "eth0",
		Destination: &TL.Pipe.LinuxBridge.UpstreamInterface,
	},

	&cli.BoolFlag{
		Name:        "linux-bridge.use-dhcp",
		Usage:       "Use the upstream DHCP server to get ip for the bridge interface.",
		Category:    category_linux_bridge,
		Required:    false,
		EnvVars:     []string{"LINUX_BRIDGE_USE_DHCP"},
		Value:       true,
		Destination: &TL.Pipe.LinuxBridge.UseDhcp,
	},

	&cli.StringFlag{
		Name:        "linux-bridge.static-ip",
		Usage:       "Use a static IP for the bridge interface.",
		Category:    category_linux_bridge,
		Required:    false,
		EnvVars:     []string{"LINUX_BRIDGE_STATIC_IP"},
		Value:       "",
		Destination: &TL.Pipe.LinuxBridge.StaticIp,
	},

	// softether

	&cli.StringFlag{
		Name:        "softether.template",
		Usage:       "Template location for the SoftEtherVPN server.",
		Category:    category_softether,
		Required:    false,
		EnvVars:     []string{"SOFTETHER_TEMPLATE"},
		Value:       "/etc/template/vpn_server.config.tmpl",
		Destination: &TL.Pipe.SoftEther.Template,
	},

	&cli.StringFlag{
		Name:        "softether.tap-interface",
		Usage:       "Interface name for SoftEther and the server to bind to as a tap device.",
		Category:    category_softether,
		Required:    false,
		EnvVars:     []string{"SOFTETHER_TAP_INTERFACE"},
		Value:       "soft",
		Destination: &TL.Pipe.SoftEther.TapInterface,
	},

	&cli.StringFlag{
		Name:        "softether.default-hub",
		Usage:       "Default hub name for SoftEtherVPN server.",
		Category:    category_softether,
		Required:    false,
		EnvVars:     []string{"SOFTETHER_DEFAULT_HUB"},
		Value:       "DEFAULT",
		Destination: &TL.Pipe.SoftEther.DefaultHub,
	},

	// server

	&cli.StringFlag{
		Name:        "server.mode",
		Usage:       `Server mode changes the behavior of the container. enum("dhcp", "bridge")`,
		Category:    category_server,
		Required:    true,
		EnvVars:     []string{"SERVER_MODE"},
		Destination: &TL.Pipe.Server.Mode,
	},

	&cli.StringFlag{
		Name:        "server.cidr-address",
		Usage:       "CIDR address of the server.",
		Category:    category_server,
		Required:    false,
		Value:       "10.0.0.0/24",
		EnvVars:     []string{"SERVER_CIDR_ADDRESS"},
		Destination: &TL.Pipe.Server.CidrAddress,
	},
}

func ProcessFlags(tl *TaskList[Pipe]) error {
	tl.Pipe.DhcpServer.ForwardingZone = tl.CliContext.StringSlice("dhcp-server.forwarding-zone")

	return nil
}
