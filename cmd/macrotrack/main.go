package main

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"macrotrack/internal/pkg/server"

	"macrotrack/internal/pkg/store"
)

// All - return true iff len of all envVars > 0
func All(envVars ...string) bool {
	for _, envVar := range envVars {
		if len(envVar) < 1 {
			return false
		}
	}
	return true
}

// track uptime
var Healthy int64

func run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel) // set logging to how

	log.WithFields(log.Fields{}).Info("starting up with these settings")
	//	log.SetOutput(os.Stderr) // reset default output

	//read env variables``
	port := os.Getenv("PORT")

	log.WithFields(log.Fields{
		"PORT": port, // default port
	}).Info("starting up with these settings")

	if !All(port) {
		log.Fatal("env arg missing")
	}

	storageType := os.Getenv("STORAGE")
	DSN := os.Getenv("DSN")

	storage := store.GetStorage(storageType, DSN)

	atomic.StoreInt64(&Healthy, time.Now().UnixNano())

	svr := &server.Server{Store: storage, Healthy: &Healthy}

	shutdownSig := svr.Start(port, "")

	<-shutdownSig

	atomic.StoreInt64(&Healthy, 0)

	//sigChannel := make(chan os.Signal, 1)
	//signal.Notify(sigChannel, os.Interrupt)
	//<-sigChannel // kill signal  ,  // force kill fuser -k apiPort/tcp
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
