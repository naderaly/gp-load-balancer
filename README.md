# HAProxy

HAProxy (High Availability Proxy) is a TCP/HTTP load balancer and proxy  server that allows to spread incoming requests across  multiple endpoints. Instead of a client connecting to a single server which processes all of the requests, the client will connect to an HAProxy instance, which will use a reverse proxy to forward the request to one of the available  endpoints, based on a load-balancing algorithm.

## Installation

The following steps were implemented on a small ubuntu EC2 instance, You we need superuser privileges. 

```bash
Sudo su -
```

Open access to port 1344:
```bash
ufw allow 1344
```

Install the HAProxy package:

```bash
apt-get update && apt-get upgrade
apt-get install haproxy
```

## Configuration

* When you configure load balancing using HAProxy, there are two types of  nodes which need to be defined: **frontend** and **backend**. The frontend is  the node by which HAProxy listens for connections. Backend nodes are  those by which HAProxy can forward requests. A third node type, the  stats node, can be used to monitor the load balancer and the other two  nodes.

* Add the following blocks of settings to the /etc/haproxy/haproxy.cfg file:

```bash
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
```

* The balance setting specifies the load-balancing strategy. In this case, the **roundrobin** strategy is used. This strategy uses each server in turn but allows for weights to be assigned to each server: servers with higher weights are  used more frequently. Other strategies include **static-rr**, which is similar to **roundrobin** but does not allow weights to be adjusted on the fly; and **leastconn**, which will forward requests to the server with the lowest number of connections.

* The **server** lines define the actual server nodes and their IP addresses, to which IP addresses will be forwarded.

* Restart the HAProxy service & check it's status.

  ```bash
  systemctl restart haproxy.service 
  systemctl status haproxy.service 
  ```

## Testing
* To confirm functionality telnet your localhost on port 1344 & press 'Enter' multi times as follow

  ```bash
  telnet localhost 1344
  ```

* It should print the following indicating server : C-ICAP/0.5.6

  ```
  Trying 127.0.0.1...
  Connected to localhost.
  Escape character is '^]'.
                  
  
  ICAP/1.0 400 Bad request
  Server: C-ICAP/0.5.6
  Connection: close
  
  Connection closed by foreign host.
  ```

  