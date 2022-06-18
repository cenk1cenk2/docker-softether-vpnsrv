package pipe

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		Set(func(t *Task[Pipe]) error {
			// set network setup
			_, network, err := net.ParseCIDR(t.Pipe.Server.CidrAddress)

			if err != nil {
				return err
			}

			t.Pipe.Ctx.Server.Network = network

			t.Log.Debugf("Network parsed: %s", t.Pipe.Ctx.Server.Network.String())

			if t.Pipe.Ctx.Server.RangeStart, err = cidr.Host(network, 1); err != nil {
				return err
			}

			if t.Pipe.Ctx.Server.RangeEnd, err = cidr.Host(network, -2); err != nil {
				return err
			}

			t.Log.Debugf(
				"Network start address: %s, Network end address: %s",
				t.Pipe.Ctx.Server.RangeStart.String(),
				t.Pipe.Ctx.Server.RangeEnd.String(),
			)

			if t.Pipe.Ctx.Health.Duration, err = time.ParseDuration(t.Pipe.Health.CheckInterval); err != nil {
				return err
			}

			// set default health address
			if t.Pipe.Health.DhcpServerAddress == "" {
				t.Lock.Lock()
				t.Pipe.Health.DhcpServerAddress = t.Pipe.Ctx.Server.RangeStart.String()
				t.Lock.Unlock()

				t.Log.Debugf(
					"Default health address for DHCP server set as default: %s",
					t.Pipe.Health.DhcpServerAddress,
				)
			}

			return nil
		})
}

func KeepAlive(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("keep-alive").
		Set(func(t *Task[Pipe]) error {
			<-t.Pipe.Terminator.Terminated

			return nil
		})

}

func CreatePostroutingRules(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("postrouting").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand("iptables").
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				Set(func(c *Command[Pipe]) error {
					c.AppendArgs(
						"-t",
						"nat",
						"-A",
						"POSTROUTING",
						"-s",
						t.Pipe.Server.CidrAddress,
						"-j",
						"MASQUERADE",
					)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			err := t.RunCommandJobAsJobParallel()

			if err != nil {
				return err
			}

			t.Log.Infof("Created postrouting rules for: %s", t.Pipe.Server.CidrAddress)

			return err
		})
}

func GenerateDhcpServerConfiguration(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("conf:dnsmasq").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Server.Mode != SERVER_MODE_DHCP
		}).
		ShouldRunBefore(func(t *Task[Pipe]) error {

			// set default gateway address
			if t.Pipe.DhcpServer.Gateway == "" && t.Pipe.DhcpServer.SendGateway {
				t.Lock.Lock()
				t.Pipe.DhcpServer.Gateway = t.Pipe.Ctx.Server.RangeStart.String()
				t.Lock.Unlock()

				t.Log.Debugf(
					"Default gateway address for DHCP server set as default: %s",
					t.Pipe.Health.DhcpServerAddress,
				)
			}

			return nil
		}).
		Set(func(t *Task[Pipe]) error {
			linkFrom := path.Join(CONF_DIR, CONF_DNSMASQ_NAME)
			linkTo := path.Join(CONF_DNSMASQ_DIR, CONF_DNSMASQ_NAME)

			if _, err := os.Stat(linkFrom); os.IsNotExist(err) {
				// generate dnsmasq configuration
				tmpl, err := template.ParseFiles(t.Pipe.DhcpServer.Template)

				if err != nil {
					return err
				}

				output := new(bytes.Buffer)

				if err := tmpl.Execute(output, DnsMasqConfigurationTemplate{
					TapInterface:      fmt.Sprintf("tap_%s", t.Pipe.SoftEther.TapInterface),
					RangeStartAddress: t.Pipe.Ctx.Server.RangeStart.String(),
					RangeEndAddress:   t.Pipe.Ctx.Server.RangeEnd.String(),
					Gateway:           t.Pipe.DhcpServer.Gateway,
					RangeNetmask:      net.IP(t.Pipe.Ctx.Server.Network.Mask).String(),
					LeaseTime:         t.Pipe.DhcpServer.Lease,
					ForwardingZone:    t.Pipe.DhcpServer.ForwardingZone.Value(),
				}); err != nil {
					return err
				}

				f, err := os.Create(linkFrom)

				if err != nil {
					return err
				}

				defer f.Close()

				if _, err = f.Write(output.Bytes()); err != nil {
					return err
				}

				t.Log.Infof("DHCP server configuration file generated: %s", linkFrom)
			} else {
				t.Log.Infof("Persistent configuration file found: %s", linkFrom)
			}

			if err := os.Remove(linkTo); err != nil {
				t.Log.Debugf(err.Error())
			}

			if err := os.Symlink(linkFrom, linkTo); err != nil {
				return err
			}

			return nil
		})
}

func GenerateSoftEtherServerConfiguration(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("conf:softether").
		Set(func(t *Task[Pipe]) error {
			linkFrom := path.Join(CONF_DIR, CONF_SOFTETHER_NAME)
			linkTo := path.Join(CONF_SOFTETHER_DIR, CONF_SOFTETHER_NAME)

			if _, err := os.Stat(linkFrom); os.IsNotExist(err) {
				// generate softether configuration
				tmpl, err := template.ParseFiles(t.Pipe.SoftEther.Template)

				if err != nil {
					return err
				}

				output := new(bytes.Buffer)

				if err := tmpl.Execute(output, SoftEtherConfigurationTemplate{
					Interface:  t.Pipe.SoftEther.TapInterface,
					DefaultHub: t.Pipe.SoftEther.DefaultHub,
				}); err != nil {
					return err
				}

				f, err := os.Create(linkFrom)

				if err != nil {
					return err
				}

				defer f.Close()

				if _, err = f.Write(output.Bytes()); err != nil {
					return err
				}

				t.Log.Warnf("SoftEtherVPN server configuration file generated: %s", linkFrom)
			} else {
				t.Log.Infof("Persistent configuration file found: %s", linkFrom)
			}

			if err := os.Remove(linkTo); err != nil {
				t.Log.Debugf(err.Error())
			}

			if err := os.Symlink(linkFrom, linkTo); err != nil {
				return err
			}

			return nil
		})
}

func CreateTapDevice(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("interface:tap").
		ShouldRunBefore(func(t *Task[Pipe]) error {
			t.Pipe.SoftEther.TapInterface = fmt.Sprintf(
				"tap_%s",
				t.Pipe.SoftEther.TapInterface,
			)

			return nil
		}).
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"ip",
				"tuntap",
				"del",
				"dev",
				t.Pipe.SoftEther.TapInterface,
				"mode",
				"tap",
			).
				SetIgnoreError(true).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			t.CreateCommand(
				"ip",
				"tuntap",
				"add",
				"dev",
				t.Pipe.SoftEther.TapInterface,
				"mode",
				"tap",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

				// give the server static ip for dnsmasq when on dhcp mode
			if t.Pipe.Server.Mode == SERVER_MODE_DHCP {
				t.CreateCommand(
					"ifconfig",
					t.Pipe.SoftEther.TapInterface,
					t.Pipe.DhcpServer.Gateway,
					"netmask",
					net.IP(t.Pipe.Ctx.Server.Network.Mask).String(),
				).
					SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
					AddSelfToTheTask()

				t.Log.Debugf(
					"Should add gateway to the tap interface: %s -> %s",
					t.Pipe.SoftEther.TapInterface,
					t.Pipe.DhcpServer.Gateway,
				)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			err := t.RunCommandJobAsJobSequence()

			t.Log.Infof(
				"Created tap adapter: %s",
				t.Pipe.SoftEther.TapInterface,
			)

			return err
		})
}

func CreateBridgeDevice(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("interface:bridge").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Server.Mode != SERVER_MODE_BRIDGE
		}).
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				"ifconfig",
				t.Pipe.LinuxBridge.BridgeInterface,
				"down",
			).
				SetIgnoreError(true).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			t.CreateCommand(
				"brctl",
				"delbr",
				t.Pipe.LinuxBridge.BridgeInterface,
			).
				SetIgnoreError(true).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			t.CreateCommand(
				"brctl",
				"addbr",
				t.Pipe.LinuxBridge.BridgeInterface,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			t.CreateCommand(
				"brctl",
				"addif",
				t.Pipe.LinuxBridge.BridgeInterface,
				t.Pipe.LinuxBridge.UpstreamInterface,
				t.Pipe.SoftEther.TapInterface,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			t.CreateCommand(
				"ip",
				"link",
				"set",
				"dev",
				t.Pipe.LinuxBridge.BridgeInterface,
				"up",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			t.CreateCommand(
				"brctl",
				"show",
				t.Pipe.LinuxBridge.BridgeInterface,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			err := t.RunCommandJobAsJobSequence()

			t.Log.Infof(
				"Created bridge adapter: %s -> %s %s",
				t.Pipe.LinuxBridge.BridgeInterface,
				t.Pipe.SoftEther.TapInterface,
				t.Pipe.LinuxBridge.UpstreamInterface,
			)

			return err
		})
}
