package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Rosalita/k8s-service/foundation/logger"
	"github.com/ardanlabs/conf/v3"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

// At build time, linker flags are used to overwrite this variable with the build reference.
var build = "develop"

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

	// Align GOMAXPROCS in the container with the number of processors available in k8s.
	opt := maxprocs.Logger(log.Infof)
	if _, err := maxprocs.Set(opt); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	defer log.Infow("shutdown")

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown // blocks until able to read from shutdown.

	return nil
}
