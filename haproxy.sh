#!/bin/bash
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi

ufw allow 1344
ufw allow 11344
apt-get install -y haproxy
cat >> /etc/haproxy/haproxy.cfg << EOF
#The frontend is the node by which HAProxy listens for connections.
frontend ICAP
bind 0.0.0.0:1344
mode tcp
default_backend icap_pool

#Backend nodes are those by which HAProxy can forward requests
backend icap_pool
balance roundrobin
mode tcp
server icap01 54.155.107.189:1344 check
server icap02 34.240.204.39:1344 check

#Haproxy monitoring Webui(optional) configuration, access it <Haproxy IP>:32700
listen stats
bind :32700
stats enable
stats uri /
stats hide-version
stats auth username:password
EOF
systemctl restart haproxy.service 
