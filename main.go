package main

import (
	"time"
	"os"
	"chatServer/internals/logger"
	"chatServer/internals/rules"
	"chatServer/internals/rules/syncdto"
	"chatServer/internals/api"
	"github.com/joho/godotenv"
	"runtime"
)

var (
	log logger.Logger
	error_std error
)

func printStats(log *logger.Logger, interval time.Duration) {
    var m runtime.MemStats

    for {
        runtime.ReadMemStats(&m)
        log.Infof(
            "[Parameters] Goroutines: %d, Alloc: %.2fMB, TotalAlloc: %.2fMB, Sys: %.2fMB, HeapObjects: %d",
            runtime.NumGoroutine(),
            float64(m.Alloc)/1024/1024,
            float64(m.TotalAlloc)/1024/1024,
            float64(m.Sys)/1024/1024,
            m.HeapObjects,
        )
				time.Sleep(interval)
    }
}

func main() {
	
	log, error_std = logger.New(time.Now().String() + ".log")
	if error_std != nil {
		panic(error_std)
	}

	err := godotenv.Load()
  if err != nil {
    log.Errorf("Error loading .env file")
  }
	
	log.Infof("Server is starting")
	clients := syncdto.NewSafeList[rules.Client]()
	channels := syncdto.NewSafeMap[rules.Channel]()

	PORT := os.Getenv("PORT")

	go api.Dealer(&log, &clients, PORT)
	log.Infof("Dealer started")
	go api.ConnectionRead(&log, &clients, &channels)
	log.Infof("Connection manager started")
	printStats(&log, 10*time.Second)
}