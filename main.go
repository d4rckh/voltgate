package main

import (
	"log"
	"net/http"
	"voltgate-proxy/config"
	"voltgate-proxy/proxy"
)

func main() {
	println("            .__   __                __          \n___  ______ |  |_/  |_  _________ _/  |_  ____  \n\\  \\/ /  _ \\|  |\\   __\\/ ___\\__  \\\\   __\\/ __ \\ \n \\   (  <_> )  |_|  | / /_/  > __ \\|  | \\  ___/ \n  \\_/ \\____/|____/__| \\___  (____  /__|  \\___  >\n                     /_____/     \\/          \\/ ")

	proxyServer := proxy.NewProxyServer()

	initialConfig := config.LoadConfig(proxyServer, "config.yaml")

	if initialConfig.ReloadConfigInterval != 0 {
		log.Println("Config reloading enabled, interval set to", initialConfig.ReloadConfigInterval, "seconds")
		go config.ReloadConfig(proxyServer, initialConfig.ReloadConfigInterval, "config.yaml")
	}

	http.HandleFunc("/", proxyServer.HandleRequest)
	log.Println("Proxy server started on", initialConfig.Address)
	log.Fatal(http.ListenAndServe(initialConfig.Address, nil))
}
