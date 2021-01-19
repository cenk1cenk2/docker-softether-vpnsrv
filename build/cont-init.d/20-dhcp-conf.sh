#!/bin/bash

echo "Creating DHCP server configuration..."
# Create DHCP configuration
if [ ! -f /cfg/dnsmasq.conf ]; then
  cp /etc/dnsmasq.conf.default /cfg/dnsmasq.conf
fi

sed -e "s/\$SRVIPSUBNET/${SRVIPSUBNET:-10.0.0}/g" -e "s/\$SRVIPNETMASK/${SRVIPNETMASK:-255.255.255.0}/g" /cfg/dnsmasq.conf >/etc/dnsmasq.conf
