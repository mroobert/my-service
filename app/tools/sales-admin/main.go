// This program performs administrative tasks for the garage sale service.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf/v2"
	"github.com/mroobert/my-service/app/tools/sales-admin/commands"
	"github.com/mroobert/my-service/foundation/logger"
	"go.uber.org/zap"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {

	// Construct the application logger.
	log, err := logger.New("ADMIN")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		if !errors.Is(err, commands.ErrHelp) {
			log.Errorw("startup", "ERROR", err)
		}
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// =========================================================================
	// Configuration

	cfg := struct {
		conf.Version
		Args conf.Args
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
		return fmt.Errorf("parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	// =========================================================================
	// Commands

	return processCommands(cfg.Args, log)
}

// processCommands handles the execution of the commands specified on
// the command line.
func processCommands(args conf.Args, log *zap.SugaredLogger) error {
	switch args.Num(0) {
	case "genkey":
		if err := commands.GenKey(); err != nil {
			return fmt.Errorf("key generation: %w", err)
		}
	default:
		fmt.Println("genkey: generate a set of private/public key files")
		return commands.ErrHelp
	}

	return nil
}
