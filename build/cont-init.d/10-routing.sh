#!/bin/bash

# Create IP tables rules
echo "Creating postrouting rules."
iptables -t nat -A POSTROUTING -s ${SRVIPSUBNET:-10.0.0}.0/24 -j MASQUERADE
