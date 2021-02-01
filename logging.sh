#!/bin/sh
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi
apt-get install -y rsyslog rsyslog-mmjsonparse rsyslog-elasticsearch rsyslog-mmutf8fix
envsubst < 10-elasticsearch.conf.tmpl > /etc/rsyslog.d/elasticsearch.conf 
systemcl restart rsyslog