This documentation shows the steps for installing & configuring haproxy to be used as web server load balancer for [glasswallsolutions.com.glasswall-icap.com](https://glasswallsolutions.com.glasswall-icap.com).

## Installation

* After connecting to the instance, you we need superuser privileges:

```bash
Sudo su -
```

* Install the HAProxy package:

```bash
apt-get update && apt-get upgrade
apt-get install haproxy
```

## Configuration

* When you configure load balancing using HAProxy, there are two types of  nodes which need to be defined: **frontend** and **backend**. The frontend is  the node by which HAProxy listens for connections. Backend nodes are  those by which HAProxy can forward requests. A third node type, the  stats node, can be used to monitor the load balancer and the other two  nodes.

* Using your favorite text editor open & add the following blocks of settings to the /etc/haproxy/haproxy.cfg file:

```bash
#The frontend is the node by which HAProxy listens for connections (http).
frontend http-glasswall
        bind *:80
        option tcplog
        mode tcp
        default_backend http-nodes  
#Backend nodes are those by which HAProxy can forward requests
backend http-nodes
        mode tcp
        balance roundrobin
        server web01 54.78.209.23:80 check

#The frontend is the node by which HAProxy listens for connections (https).
frontend https-glasswall
        bind *:443
        option tcplog
        mode tcp
        default_backend https-nodes 
#Backend nodes are those by which HAProxy can forward requests
backend https-nodes
        mode tcp
        balance roundrobin
        option ssl-hello-chk
        server web01 54.78.209.23:443 check

#Haproxy monitoring Webui(optional) configuration, access it <Haproxy IP>:32700
listen stats
bind :32700
option http-use-htx
http-request use-service prometheus-exporter if { path /metrics }
stats enable
stats uri /
stats hide-version
stats auth username:password
```

* Here we have two frontends as HAproxy is listening on port 80 for http connection and on port 443 for secure https connection

* The balance setting specifies the load-balancing strategy. In this case, the **roundrobin** strategy is used. This strategy uses each server in turn but allows for weights to be assigned to each server: servers with higher weights are  used more frequently. Other strategies include **static-rr**, which is similar to **roundrobin** but does not allow weights to be adjusted on the fly; and **leastconn**, which will forward requests to the server with the lowest number of connections.

* The **server** lines define the actual server nodes and their IP addresses, to which IP addresses will be forwarded.

* The stats section defines HAProxy **Stats page** that shows you an abundance of metrics that cover the health of your servers, current request rates, response times, and more.
  The HAProxy stats node will listen on port 32700 for connections and is  configured to hide the version of HAProxy as well as to require a  password login. Replace `password` with a more secure password.
  
  From your browser access the stats page through the following:
  URL : <Haproxy IP>:32700
  
  Username: username
  
  Password: password 

![image](https://user-images.githubusercontent.com/58347752/101278958-7cfe6b00-37c7-11eb-923d-28788f224433.png)

* Restart the HAProxy service & check it's status.

  ```bash
  systemctl restart haproxy.service 
  systemctl status haproxy.service 
  ```

## Client configuration

- Add hosts records to your client system hosts file ( i.e **Windows**: C:\Windows\System32\drivers\etc\hosts , **Linux, macOS and  Unix-like:** /etc/hosts ) as follows

```
ip a

<VM IP ADDRESS> glasswallsolutions.com.glasswall-icap.com
```

make sure that tcp ports **80** and **443** are reachable and not blocked by firewall.

## Access the proxied site

* You can access the proxied site by browsing [glasswallsolutions.com.glasswall-icap.com](https://glasswallsolutions.com.glasswall-icap.com).
* Verify that your access is established through the haproxy loadbalancer through the network tab, the Request address should show the Loadbalancer server IP as shown below

![image](https://user-images.githubusercontent.com/58347752/100607205-4500af00-3313-11eb-8f14-b075e74108a7.png)

