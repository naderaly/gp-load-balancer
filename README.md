# HAProxy

HAProxy (High Availability Proxy) is a TCP/HTTP load balancer and proxy  server that allows to spread incoming requests across multiple endpoints. Instead of a client connecting to a single server which processes all of the requests, the client will connect to an HAProxy instance, which will use a reverse proxy to forward the request to one of the available endpoints, based on a configured load-balancing algorithm.

Please refer to [haproxy-icap.md](https://github.com/k8-proxy/gp-load-balancer/blob/main/haproxy-icap.md) for the steps of installing & configuring haproxy to be used as ICAP server load balancer.

Or refer to [haproxy-web.md](https://github.com/k8-proxy/gp-load-balancer/blob/main/haproxy-web.md) for the steps of installing & configuring haproxy to be used as web server load balancer.

# **Health check**

check the server by sending test.pdf file to server 

if the response with HTTP/1.0 200 OK so this will be healthy server and exitcode will be zero

and will return reb_test.pdf  file

else not response or forbidden so this will be not  healthy server and exitcode will be one

time out for command 32 second 

[**Requirement]()**

> #### to run command you need to have c-icap-client 

[**Running the tests]()**

go run check.go 192.168.100.100 80 54.194.133.136 1344 -v

```
the third argument 54.194.133.136 server ip fourth argument port number the last  argument -v verbos out of shell script  
```

```html
  <img src="./image/healthy.jpg" width="350" title="healthy server">
```

exitcode 0 it's  successful and healthy server 

go run check.go 192.168.100.100 80 34.244.7.158 433 -v
```html
  <img src="./image/error.jpg" width="350" title="healthy server">
```

exitcode 1 it's failed and not healthy server

