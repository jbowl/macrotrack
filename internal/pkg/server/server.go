package server

import (
	"context"
	"log"

	"macrotrack/internal/pkg/handlers"
	"macrotrack/internal/pkg/store"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	//	bc      *hodlapi.BreweryServiceClient
	Store   store.Storage
	Healthy *int64
}

// Start - start http server and grpc client , returns a signal channel
func (s *Server) Start(port string, apiAddr string) <-chan os.Signal {
	//// grpc
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	err := s.Store.Init()

	if err != nil {
		log.Fatal(err)
	}
	go func() {
		// brpc
		router := mux.NewRouter().StrictSlash(true)

		//router.HandleFunc()

		// Create
		router.Handle("/macros", handlers.CreateMacro(s.Store)).Methods("POST")
		// Read

		router.Handle("/macros", handlers.ReadAllMacro(s.Store)).Methods("GET") // paginate quantity

		router.Handle("/macros/{uuid}", handlers.ReadMacro(s.Store)).Methods("GET")
		//
		// Update
		router.Handle("/macros/{uuid}", handlers.UpdateMacro(s.Store)).Methods("PUT")
		// Delete
		router.Handle("/macros/{uuid}", handlers.DeleteMacro(s.Store)).Methods("DELETE")

		//router.Handle("/breweries/search", handlers.Search()).Methods("GET")
		//router.Handle("/healthz", handlers.HealthZ(s.Healthy))

		//  potential alternative to the heavy approach of httputil.NewSingleHostReverseProxy

		// forwarded to another service
		//router.PathPrefix("/mobile").HandlerFunc(mobileReq).Methods("GET", "POST", "PUT", "PATCH", "OPTIONS", "DELETE")

		router.NotFoundHandler = handlers.Custom404Handler()
		router.MethodNotAllowedHandler = handlers.Custom405Handler()

		router.Use(mux.CORSMethodMiddleware(router))

		httpServer := http.Server{
			Addr:         ":" + port,
			ReadTimeout:  5 * time.Minute,
			WriteTimeout: 5 * time.Minute,
			IdleTimeout:  5 * time.Minute,
			Handler:      router,
		}

		log.Fatal(httpServer.ListenAndServe())
		//	wg.Done()
	}()

	return shutdown
}

func (s *Server) Shutdown() {
	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		// Release resources like Database connections
		cancel()
	}()

	shutdownChan := make(chan error, 1)
	// TODO :
	//	go func() { shutdownChan <- s.App.Shutdown() }()

	select {
	case <-timeout.Done():
		log.Fatal("Server Shutdown Timed out before shutdown.")
	case err := <-shutdownChan:
		if err != nil {
			log.Fatal("Error while shutting down server", err)
		} else {
			log.Println("Server Shutdown Successful")
		}
	}
}
