#!/bin/bash

source /scripts/logger.sh

log_start "Creating SoftEther server configuration..."

# Remove the default configuration
rm /etc/vpnserver/vpn_server.config 2>/dev/null

# Create default configuration
if [ ! -f /cfg/vpn_server.config ]; then
  log_warn "Using default SoftEther configuration, since no persistent configuration found."

  cp /default/vpnserver/vpn_server.config /cfg/vpn_server.config
fi

# Link files
ln -s /cfg/vpn_server.config /etc/vpnserver
