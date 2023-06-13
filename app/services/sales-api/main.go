package main

import (
	"fmt"
	"os"

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
	return nil
}
