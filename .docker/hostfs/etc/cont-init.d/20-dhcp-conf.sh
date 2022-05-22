#!/bin/bash

source /scripts/logger.sh

log_start "Creating DHCP server configuration..."

# Create DHCP configuration
if [ ! -f /cfg/dnsmasq.conf ]; then
  log_warn "Using default DNSMASQ configuration, since no persistent configuration found."

  cp /default/dnsmasq/dnsmasq.conf /cfg/dnsmasq.conf
fi

# Swapping out variables from the DNSMASQ config file
sed \
  -e "s/\$SRVIPSUBNET/${SRVIPSUBNET:-10.0.0}/g" \
  -e "s/\$SRVIPNETMASK/${SRVIPNETMASK:-255.255.255.0}/g" \
  -e "s/\$DHCP_START/${DHCP_START:-10}/g" \
  -e "s/\$DHCP_END/${DHCP_END:-254}/g" \
  -e "s/\$DHCP_LEASE/${DHCP_END:-"12h"}/g" \
  /cfg/dnsmasq.conf >/etc/dnsmasq.conf

log_finish "Injected server variables for DNSMASQ configuration."
