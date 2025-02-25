package main

import (
	"log"
	"net/http"
	"voltgate-proxy/config"
	"voltgate-proxy/proxy"
)

func main() {
	println("            .__   __                __          \n___  ______ |  |_/  |_  _________ _/  |_  ____  \n\\  \\/ /  _ \\|  |\\   __\\/ ___\\__  \\\\   __\\/ __ \\ \n \\   (  <_> )  |_|  | / /_/  > __ \\|  | \\  ___/ \n  \\_/ \\____/|____/__| \\___  (____  /__|  \\___  >\n                     /_____/     \\/          \\/ ")

	// Load configuration
	appConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Create proxy server instance
	proxyServer := proxy.NewProxyServer(appConfig)

	// Start HTTP server
	http.HandleFunc("/", proxyServer.HandleRequest)
	log.Println("Proxy server started on :80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

// Handle incoming requests and forward to mapped services
