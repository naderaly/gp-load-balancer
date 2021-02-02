# **HAProxy External Health check**

HAProxy minimum configuration for a health check is the `check` keyword on a `server` line. In order to run, a health check requires at least an **IP address** and a **TCP port** from the server.

Though a TCP test connection is enough in some cases but it only test the health of the server and doesn't test the health a service running on that server.

So HAProxy present external health check option which basically an option which HAProxy  runs an external script and according  to it's exit code it decides whether the server is alive and continue sending traffic to or the server is considered unhealthy and stop sending  traffic to it.

If the script's exit code is 0 the server is healthy, if not HAProxy mark the server as unhealthy.

## Configuring HAProxy to run  external health check for ICAP servers

First, we need to set some prerequisites 

* Install **golang** on your server

  ```bash
  sudo wget https://golang.org/dl/go1.15.7.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.15.7.linux-amd64.tar.gz
  echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
  export PATH=$PATH:/usr/local/go/bin
  ```

* Clone the repository and build the script

  ```bash
  git clone https://github.com/naderaly/gp-load-balancer
  cd gp-load-balancer/External-HealthCheck
  go build check.go
  sudo cp check test.pdf /var/lib/haproxy/dev
  sudo chown haproxy:haproxy /var/lib/haproxy/dev
  ```

* Open HAProxy configuration file 

  ```bash
  sudo vim /etc/haproxy/haproxy.cfg
  ```

  * Delete or comment out the following line

    ```bash
    chroot /var/lib/haproxy
    ```

  * Add the following two lines to the backend you want to apply external healthcheck to

    ```bash
    option external-check
    external-check command /var/lib/haproxy/dev/check.sh
    ```

  * for example:

    ```bash
    #Backend nodes are those by which HAProxy can forward requests
    backend icap_pool
    balance roundrobin
    mode tcp
    option external-check
    external-check command /var/lib/haproxy/dev/check.sh
    server icap01-es 51.89.210.148:1344 check inter 60s fall 3 rise 2
    server icap02-es 51.89.210.149:1344 check inter 60s
    server icap03-es 51.89.210.151:1344 check inter 60s
    server icap04-es 51.89.210.152:1344 check inter 60s
    server icap05-es 51.89.210.155:1344 check inter 60s
    ```

  * Health check parameters:

    All the keywords below apply to the `server` or `default-server` directives which you can use to set up a health check frequency based on the server's state:

    ```tex
    inter        : Sets the interval between two consecutive health checks. If not specified, the default value is 2s.
    rise <count> : Number of consecutive valid health checks before considering the server as UP. Default value is 2
    fall <count> : Number of consecutive invalid health checks before considering the server as DOWN. Default value is 
    ```

* Save the file and quite and reload HAProxy

  ```bash
  sudo systemctl reload haproxy
  ```

  