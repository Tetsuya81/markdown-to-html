package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 8080, "Port to serve on")
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}

	fs := http.FileServer(http.Dir(dir))

	http.Handle("/", fs)

	serverAddr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting server at http://localhost%s\n", serverAddr)
	fmt.Printf("Serving files from: %s\n", dir)
	fmt.Println("Press Ctrl+C to stop the server")
	
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
