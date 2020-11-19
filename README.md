# HAProxy

HAProxy (High Availability Proxy) is a TCP/HTTP load balancer and proxy  server that allows to spread incoming requests across multiple endpoints. Instead of a client connecting to a single server which processes all of the requests, the client will connect to an HAProxy instance, which will use a reverse proxy to forward the request to one of the available endpoints, based on a configured load-balancing algorithm.

Please refer to [haproxy-icap.md](https://github.com/k8-proxy/gp-load-balancer/blob/main/haproxy-icap.md) for the steps of installing & configuring haproxy to be used as ICAP server load balancer.

