# Round-Robin-Load-Balancer
A simple round robin load balancer implemented in Go

Here we supply a list of routes to servers that generate a load balancer that uses round robin selection to route requests to each server.

**Example**

1. go run any number of servers from /server
2. go run loadbalancer/loadbalancer.go

Now when you head over to localhost:8080. Based on the servers started it should show a new port every time you refresh the page.

**Methodology**

It works by looping through all available servers using round robin methodology and then routes the request to the appropriate server.

In a way using the load balancer here is fault tolerant. In the case a server goes down, the load balancer routes the request to the next available server. Logs are shown in loadbalancer.go to show the routing mechanisms.

Health checks are done in real time during every route to determine which servers are alive or not. Only in the case when all the given servers are down, it returns a http.StatusServiceUnavailable status back to the client.

**Future Improvements**

* Generate statistics for each server attached and then route accordingly using least connections method.
* Auto scale the number of servers based on the above statistics.
* Use goroutines to continuously do health checks on each server. 
* A CRUD based interface to add, remove and update servers in the load balancer.

