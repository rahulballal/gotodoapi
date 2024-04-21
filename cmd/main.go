package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rahulballal/gotodoapi/internal"
)

func main() {
	configPtr := internal.LoadConfig()
	mux := http.NewServeMux()
	logger := log.Default()
	internal.InitializePersistence()
	handlerMap := &internal.HandlerMap{}
	handlerMap.Logger = logger
	internal.ConfigureRouting(mux, handlerMap)
	address := fmt.Sprintf(":%d", configPtr.Port)
	serverInitError := http.ListenAndServe(address, mux)
	log.Fatal(serverInitError)
}
