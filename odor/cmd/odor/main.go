package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jlorgal/odor/odor"
	"github.com/jlorgal/odor/odor/filters"
	"github.com/jlorgal/odor/odor/profile"
	"github.com/jlorgal/odor/odor/svc"
)

func main() {
	// Prepare logger
	time.Local = time.UTC
	logContext := svc.LogContext{
		Service:   "odor",
		Operation: "init",
	}
	logger := svc.NewLogger()
	logger.SetLogContext(&logContext)

	// Prepare the configuration
	cfgFile := flag.String("config", "./config.json", "path to config file")
	flag.Parse()
	var cfg odor.Config
	if err := svc.GetConfig(*cfgFile, &cfg); err != nil {
		logger.Fatal("Bad configuration with file '%s'. %s", *cfgFile, err)
		os.Exit(1)
	}
	logger.SetLevel(cfg.LogLevel)

	// Log the configuration
	if configBytes, err := json.Marshal(cfg); err == nil {
		logger.Info("Configuration: %s", string(configBytes))
	}

	serviceProfile := profile.New(&cfg)

	// Capture signals to stop the service
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		for sig := range c {
			logger.Warn("Captured signal %s. Stopping service", sig)
			// TODO: How to stop the service???
			serviceProfile.Stop()
			os.Exit(0)
		}
	}()

	// Start the Profile Service
	logger.Info("Starting service")
	serviceProfile.Start()

	// Start the pipeline engine

	pipeline := odor.NewPipeline()
	pipeline.AddFilters(
		filters.NewParentalControl(),
	)

}
