#!/bin/bash
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi
apt update
apt upgrade -y
ufw allow 1344
ufw allow 1345
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
server icap01 54.77.168.168:1344 check
server icap02 3.139.22.215:1344 check

#The frontend is the node by which HAProxy listens for connections.
frontend S-ICAP
bind 0.0.0.0:1345
mode tcp
default_backend s-icap_pool
#Backend nodes are those by which HAProxy can forward requests
backend s-icap_pool
balance roundrobin
mode tcp
server icap01 54.77.168.168:1345 check
server icap02 3.139.22.215:1345 check

#Haproxy monitoring Webui(optional) configuration, access it <Haproxy IP>:32700
listen stats
bind :32700
stats enable
stats uri /
stats hide-version
stats auth username:password
EOF
systemctl restart haproxy.service 
