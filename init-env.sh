#!/bin/bash

echo "init-env.sh@bash: v2.0, 20190321"
## Variables
# Write down the data required in .env file here for initiation.
ENVFILENAME=.env
ENVFILECONTENTS=(
  "# Timezone"
  "TZ="
  "# Sleep Time for Server Alive Check in Seconds (default: 600)"
  "SLEEPTIME="
  "# Keep logs or delete them in between sleeptime. To keep set the type to 1."
  "KEEP_SERVER_LOG="
  "KEEP_PACKET_LOG="
  "KEEP_SECURITY_LOG="
  "# VPN Server IP Subnet in form of xx.xx.xx (default: 10.0.0), it can also rewrite dnsmasq.conf with SED if \$SRVIPSUBNET inside dnsmasq.conf is set."
  "SRVIPSUBNET="
  "# VPN Server IP Subnet Netmask in form of xx.xx.xx.xx (default: 255.255.255.0) \$SRVIPNETMASK"
  "SRVIPNETMASK="
  "# IP lease start address (default: 10)"
  "DHCP_START="
  "# IP lease end address (default: 254)"
  "DHCP_END="
  "# IP lease time (default: 12h)"
  "DHCP_LEASE="
)

## Script
echo "Initiating ${ENVFILENAME} file."
if [[ ! -f ${ENVFILENAME} ]] || (echo -n ".env file already initiated. You want to override? [ y/N ]: " && read -r OVERRIDE && echo ${OVERRIDE::1} | grep -iqF "y"); then
  echo "Will rewrite the .env file with the default one."
  >${ENVFILENAME} && for i in "${ENVFILECONTENTS[@]}"; do echo $i >>${ENVFILENAME}; done
  echo "Opening enviroment file in nano editor."
  nano ${ENVFILENAME}
  echo "All done."
else
  echo "File already exists with no overwrite permission given."
  echo "Not doing anything."
fi
