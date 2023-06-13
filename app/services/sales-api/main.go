package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Rosalita/kind/foundation/logger"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.New("SALES_API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	// =========================================================================
	// GOMAXPROCS
	//
	// GOMAXPROCS gets or sets the number of operating system threads the go program will use
	// to execute Go routines. Passing 0 reports the maximum number of CPUs that can be 
	// executing simultaneously. CPU is at 100% capacity when it is runnning the 
	// GOMAXPROCS number of Go routines in parallel. The number of logical processors and cores
	// should always be the same to avoid cycles being wasted by context switching.

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	defer log.Infow("shutdown")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown // blocks until able to read from shutdown.

	return nil
}
