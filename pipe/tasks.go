package pipe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/sirupsen/logrus"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type Ctx struct {
	Health struct {
		Duration time.Duration
	}

	Server struct {
		Network    *net.IPNet
		RangeStart net.IP
		RangeEnd   net.IP
	}

	DhcpServer struct {
		Options map[string]string
	}
}

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

			return nil
		})
}

func CreatePostroutingRules(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("postrouting").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand("iptables").
				SetLogLevel(logrus.DebugLevel, 0, logrus.DebugLevel).
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
	return tl.CreateTask("dnsmasq-conf").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Server.Mode != SERVER_MODE_DHCP
		}).
		ShouldRunBefore(func(t *Task[Pipe]) error {
			// set default health address
			if t.Pipe.Health.DhcpServerAddress == "" {
				t.Lock.Lock()
				t.Pipe.Health.DhcpServerAddress = t.Pipe.Ctx.Server.RangeStart.String()
				t.Lock.Unlock()

				t.Log.Infof(
					"Default health address for DHCP server set as default: %s",
					t.Pipe.Health.DhcpServerAddress,
				)
			}

			// set default gateway address
			if t.Pipe.DhcpServer.Gateway == "" && t.Pipe.DhcpServer.SendGateway {
				t.Lock.Lock()
				t.Pipe.DhcpServer.Gateway = t.Pipe.Ctx.Server.RangeStart.String()
				t.Lock.Unlock()

				t.Log.Infof(
					"Default gateway address for DHCP server set as default: %s",
					t.Pipe.Health.DhcpServerAddress,
				)
			}

			// unmarshal dhcp server options from json
			t.Lock.Lock()
			err := json.Unmarshal([]byte(t.Pipe.DhcpServer.Options), &t.Pipe.Ctx.DhcpServer.Options)
			t.Lock.Unlock()

			if err != nil {
				return err
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
					TapInterface:      fmt.Sprintf("tap_%s", t.Pipe.DhcpServer.TapInterface),
					RangeStartAddress: t.Pipe.Ctx.Server.RangeStart.String(),
					RangeEndAddress:   t.Pipe.Ctx.Server.RangeEnd.String(),
					Gateway:           t.Pipe.DhcpServer.Gateway,
					RangeNetmask:      net.IP(t.Pipe.Ctx.Server.Network.Mask).String(),
					LeaseTime:         t.Pipe.DhcpServer.Lease,
					ForwardingZone:    t.Pipe.DhcpServer.ForwardingZone.Value(),
					Options:           t.Pipe.Ctx.DhcpServer.Options,
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
				return err
			}

			if err := os.Symlink(linkFrom, linkTo); err != nil {
				return err
			}

			return nil
		})
}

func CreateTapDevice(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("interface-tap").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Server.Mode != SERVER_MODE_DHCP
		}).
		ShouldRunBefore(func(t *Task[Pipe]) error {
			t.Pipe.DhcpServer.TapInterface = fmt.Sprintf("tap_%s", t.Pipe.DhcpServer.TapInterface)

			return nil
		}).
		Set(func(t *Task[Pipe]) error {

			t.CreateCommand(
				"ip",
				"tuntap",
				"add",
				"dev",
				t.Pipe.DhcpServer.TapInterface,
				"mode",
				"tap",
			).
				SetLogLevel(logrus.DebugLevel, 0, logrus.DebugLevel).
				AddSelfToTheTask()

			t.CreateCommand(
				"ifconfig",
				t.Pipe.DhcpServer.TapInterface,
				t.Pipe.DhcpServer.Gateway,
				"netmask",
				net.IP(t.Pipe.Ctx.Server.Network.Mask).String(),
			).
				SetLogLevel(logrus.DebugLevel, 0, logrus.DebugLevel).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			err := t.RunCommandJobAsJobSequence()

			t.Log.Infof(
				"Created tap adapter: %s -> %s",
				t.Pipe.DhcpServer.TapInterface,
				t.Pipe.DhcpServer.Gateway,
			)

			return err
		})
}
