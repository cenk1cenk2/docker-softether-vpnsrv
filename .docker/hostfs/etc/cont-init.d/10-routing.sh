#!/bin/bash

source /scripts/logger.sh

log_start "Creating postrouting rules."

# Create IP tables rules
function convMask() {
  c=0
  x=0$(printf '%o' ${1//./ })

  while [ $x -gt 0 ]; do
    let c+=$((x % 2)) 'x>>=1'
  done
}

CIDR_NET=$(convMask "${SRVIPNETMASK:-255.255.255.0}")

iptables -t nat -A POSTROUTING -s ${SRVIPSUBNET:-10.0.0}.0/${CIDR_NET:-24} -j MASQUERADE

log_finish "Created postrouting for: ${SRVIPSUBNET:-10.0.0}.0/${CIDR_NET:-24}"
