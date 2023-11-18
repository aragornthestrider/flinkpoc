package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"transactiongenerator/config"
	"transactiongenerator/server"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger, err := NewLogger("info")
	if err != nil {
		log.Fatal("Error in setting up logger")
		return
	}

	go func() {
		exitSignal := <-exit
		logger.Info("Service shutting down main due to: " + exitSignal.String())
		cancel()
	}()

	configPath := flag.String("config", "/config/config.yaml", "Path of config file")
	config := config.ParseConfig(*configPath, logger)
	logger.Info("Config read is: ", zap.Any("Config", config))

	server := server.NewServer(ctx, 8080, 10, 10, logger)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		err = server.Run()
		if err != nil {
			logger.Error("Error in running server", zap.Error(err))
		}
	}()

	server.SetHealthy(true)
	server.SetReady(true)

	wg.Wait()
}

func NewLogger(configLogLevel string) (*zap.Logger, error) {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	logLevel := parseLogLevel(configLogLevel)
	core := ecszap.NewCore(encoderConfig, os.Stdout, logLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger, nil
}

func parseLogLevel(logLevel string) zapcore.Level {
	switch logLevel {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
