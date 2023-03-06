package pipe

var ()

const (
	SERVER_MODE_DHCP   = "dhcp"
	SERVER_MODE_BRIDGE = "bridge"

	CONF_DIR = "/conf/"

	CONF_DNSMASQ_NAME = "dnsmasq.conf"
	CONF_DNSMASQ_DIR  = "/etc/"

	CONF_SOFTETHER_NAME = "vpn_server.config"
	CONF_SOFTETHER_DIR  = "/etc/softether"

	HOOKS_DIR            = "/docker.init.d"
	HOOK_POST_TASKS_FILE = "post-tasks"
)
