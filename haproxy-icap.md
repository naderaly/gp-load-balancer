This documentation shows the steps for installing & configuring haproxy to be used as ICAP server load balancer for icap01.glasswall-icap.com & icap02.glasswall-icap.com. 

## Auto Installation

A script to auto install HAProxy to be used as load balancer for ICAP servers with the default configuration (as mentioned bellow) to be installed on a running ubuntu server in included in this [repository](https://github.com/k8-proxy/gp-load-balancer) 

Clone the repository and run the script as root.

```bash
cd
git clone https://github.com/k8-proxy/gp-load-balancer
sudo su -
#Please replace the placeholder with your username
cd /home/<username>/gp-load-balancer
./haproxy.sh
```

Make sure haproxy is active and running 

```bash
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


Or for manual installation or for extra configuration manipulation other than the default ones, please follow the following steps:

## Manual Installation

* The following steps were implemented on a small ubuntu 20.04 EC2 instance, We need to add an inboud rule for TCP traffic on ports 1344 &11344. 

* In the navigation pane of the Amazon EC2 console, choose **Instances**, Select your instance and look at the **Description** tab, **Security groups** lists the security groups that are associated with the instance. Choose **view inbound rules** to display a list of the rules that are in effect for the instance.

![image](https://user-images.githubusercontent.com/58347752/98373169-86b97500-2047-11eb-9115-459ffa9a08a7.png)

* Select the security group and edit inboud rules

![image](https://user-images.githubusercontent.com/58347752/98373855-85d51300-2048-11eb-91ec-03baf8568d96.png)

* Now back to your Instance page and under the Instance summery tab click on connect

![image](https://user-images.githubusercontent.com/58347752/98374258-13186780-2049-11eb-97a7-0cacc64f06ca.png)

* Follow the steps under SSH client tab to connect to the machine

![image](https://user-images.githubusercontent.com/58347752/98374350-3216f980-2049-11eb-8062-cc65dd3841fc.png)

* the .pem file should be downloaded to you localhost while creating the instance but you can redownload it by opening the Amazon EC2 console at                                                      https://console.aws.amazon.com/ec2/, Then In the navigation pane, under **NETWORK & SECURITY**, choose **Key Pairs**.                                                                                        

  ![image](https://user-images.githubusercontent.com/58347752/98374826-eadd3880-2049-11eb-9ae2-cf560df9f32b.png)

* After connectiong to the instance, you we need superuser privileges:

```bash
Sudo su -
```

* Open access to ports 1344 & 11344:

```bash
ufw allow 1344
ufw allow 11344
```

* Install the HAProxy package:

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

  