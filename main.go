package main

import (
	"time"
	"chatServer/internals/logger"
)

var (
	log logger.Logger
	error_std error
)

func init() {
	log, error_std = logger.New(time.Now().String() + ".log")
	if error_std != nil {
		panic(error_std)
	}
	
	log.Infof("Server is starting")
}

func runApp() {
	
}

func main() {
	runApp()
}