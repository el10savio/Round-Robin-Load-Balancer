package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Server struct {
	Route        string
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

type ServerList struct {
	Servers []Server
	Latest  int
}

func (server *Server) isAlive() bool {
	timeout := time.Duration(1 * time.Second)

	log.Println("Started Health Check For:", server.Route)
	_, err := net.DialTimeout("tcp", server.Route, timeout)
	if err != nil {
		log.Println(server.Route, "Is Dead")
		log.Println("Health Check Error:", err)
		server.Alive = false
		return false
	}

	log.Println(server.Route, "Is Alive")
	server.Alive = true
	return true
}

func (serverList *ServerList) Init(serverRoutes []string) {
	log.Println("Creating Server List For Routes:", serverRoutes)

	for _, serverRoute := range serverRoutes {
		var localServer Server

		localServer.Route = serverRoute
		localServer.Alive = localServer.isAlive()

		origin, _ := url.Parse("http://" + serverRoute)
		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)
			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}
		localServer.ReverseProxy = &httputil.ReverseProxy{Director: director}

		log.Println("Server", localServer, "Added To Server List")
		serverList.Servers = append(serverList.Servers, localServer)
	}

	serverList.Latest = -1
	log.Println("Successfully Created Server List:", serverList)

}

func (serverList *ServerList) nextServer() int {
	return (serverList.Latest + 1) % len(serverList.Servers)
}

func (serverList *ServerList) loadBalance(w http.ResponseWriter, r *http.Request) {
	if len(serverList.Servers) > 0 {
		serverCount := 0
		for index := serverList.nextServer(); serverCount < len(serverList.Servers); index = serverList.nextServer() {
			if serverList.Servers[index].isAlive() {
				log.Println("Routing Request", r.URL, "To", serverList.Servers[index].Route)
				serverList.Servers[index].ReverseProxy.ServeHTTP(w, r)

				serverList.Latest = index
				log.Println("Updated Latest Server To:", serverList.Latest)

				return
			}
			serverCount++
			serverList.Latest = serverList.nextServer()
		}
	}
	log.Println("No Servers Available")
	http.Error(w, "No Servers Available", http.StatusServiceUnavailable)
}

func main() {
	var serverList ServerList
	loadBalancerPort := "8080"

	serverRoutes := []string{
		"localhost:8081",
		"localhost:8083",
		"localhost:8085",
	}

	serverList.Init(serverRoutes)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serverList.loadBalance(w, r)
	})

	http.ListenAndServe(":"+loadBalancerPort, nil)
}
