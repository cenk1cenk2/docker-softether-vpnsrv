#!/bin/bash

rm /etc/vpnserver/vpn_server.config 2>/dev/null

# Create default configuration
if [ ! -f /cfg/vpn_server.config ]; then
  echo "Creating default configuration."
  cp /etc/vpnserver/vpn_server.config.default /cfg/vpn_server.config
fi

# Link files
ln -s /cfg/vpn_server.config /etc/vpnserver
