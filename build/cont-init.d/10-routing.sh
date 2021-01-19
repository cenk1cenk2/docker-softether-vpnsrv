#!/bin/bash

# Create IP tables rules
function convMask () {
  c=0 x=0$( printf '%o' ${1//./ } )
  while [ $x -gt 0 ]; do
      let c+=$((x%2)) 'x>>=1'
  done
  echo "$c";
}
echo "Creating postrouting rules."
cidrNet=$(convMask ${SRVIPNETMASK:-255.255.255.0})
iptables -t nat -A POSTROUTING -s ${SRVIPSUBNET:-10.0.0}.0/${cidrNet:-24} -j MASQUERADE
